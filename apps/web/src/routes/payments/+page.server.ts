import { fetchAPI, sendAPI } from '$lib/server/api';
import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

type Charge = {
	id: string;
	householdId: string;
	description: string;
	amountMinor: number;
	allocatedMinor: number;
	outstandingMinor: number;
	currency: string;
	status: string;
	dueDate?: string | null;
};

type Payment = {
	id: string;
	payerContactId?: string | null;
	receivedAt: string;
	amountMinor: number;
	currency: string;
	method: string;
	reference?: string | null;
	notes?: string | null;
};

type Household = { id: string; displayName: string };
type Contact = { id: string; displayName: string };

export const load: PageServerLoad = async ({ fetch }) => {
	const [charges, payments, households, contacts] = await Promise.all([
		fetchAPI<{ charges: Charge[] }>(fetch, '/api/charges?status=open', { charges: [] }),
		fetchAPI<{ payments: Payment[] }>(fetch, '/api/payments', { payments: [] }),
		fetchAPI<{ households: Household[] }>(fetch, '/api/households', { households: [] }),
		fetchAPI<{ contacts: Contact[] }>(fetch, '/api/contacts', { contacts: [] })
	]);

	return {
		apiConnected:
			charges.connected && payments.connected && households.connected && contacts.connected,
		apiError: charges.error ?? payments.error ?? households.error ?? contacts.error,
		charges: charges.data.charges,
		payments: payments.data.payments,
		householdNames: Object.fromEntries(
			households.data.households.map((household) => [household.id, household.displayName])
		),
		contactNames: Object.fromEntries(
			contacts.data.contacts.map((contact) => [contact.id, contact.displayName])
		),
		contacts: contacts.data.contacts
	};
};

const positiveAmountMinor = (value: FormDataEntryValue | null): number | null => {
	if (typeof value !== 'string' || !/^\d+(?:\.\d{1,2})?$/.test(value.trim())) return null;
	const amountMinor = Math.round(Number(value) * 100);
	return Number.isSafeInteger(amountMinor) && amountMinor > 0 ? amountMinor : null;
};

export const actions: Actions = {
	recordPayment: async ({ fetch, request }) => {
		const formData = await request.formData();
		const amountMinor = positiveAmountMinor(formData.get('amount'));
		const chargeIds = formData.getAll('chargeId').map(String);
		if (amountMinor === null) return fail(400, { error: 'Enter a positive payment amount.' });
		if (new Set(chargeIds).size !== chargeIds.length) {
			return fail(400, { error: 'A charge can only be selected once.' });
		}

		const chargeResult = await fetchAPI<{ charges: Charge[] }>(fetch, '/api/charges?status=open', {
			charges: []
		});
		if (!chargeResult.connected) return fail(503, { error: 'The API is unavailable.' });
		const chargeById = new Map(chargeResult.data.charges.map((charge) => [charge.id, charge]));
		const allocations: Array<{ chargeId: string; amountMinor: number }> = [];
		let currency = 'MXN';
		for (const chargeId of chargeIds) {
			const charge = chargeById.get(chargeId);
			const allocationMinor = positiveAmountMinor(formData.get(`allocation:${chargeId}`));
			if (!charge || allocationMinor === null || allocationMinor > charge.outstandingMinor) {
				return fail(400, { error: 'Each allocation must fit the selected charge balance.' });
			}
			if (allocations.length > 0 && charge.currency !== currency) {
				return fail(400, { error: 'All selected charges must use the same currency.' });
			}
			currency = charge.currency;
			allocations.push({ chargeId, amountMinor: allocationMinor });
		}
		const allocatedMinor = allocations.reduce(
			(total, allocation) => total + allocation.amountMinor,
			0
		);
		if (allocatedMinor > amountMinor) {
			return fail(400, { error: 'Allocations cannot exceed the payment amount.' });
		}

		const receivedAtInput = String(formData.get('receivedAt') ?? '');
		const receivedAtDate = receivedAtInput ? new Date(receivedAtInput) : new Date();
		if (Number.isNaN(receivedAtDate.getTime())) {
			return fail(400, { error: 'Choose a valid payment date and time.' });
		}
		const receivedAt = receivedAtDate.toISOString();
		const result = await sendAPI<Payment>(fetch, '/api/payments', 'POST', {
			payerContactId: String(formData.get('payerContactId') ?? '') || undefined,
			receivedAt,
			amountMinor,
			currency,
			method: String(formData.get('method') ?? 'cash'),
			reference: String(formData.get('reference') ?? '').trim() || undefined,
			notes: String(formData.get('notes') ?? '').trim() || undefined,
			allocations
		});
		if (result.error) return fail(400, { error: result.error });
		throw redirect(303, '/payments');
	}
};
