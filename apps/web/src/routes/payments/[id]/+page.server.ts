import { fetchAPI, sendAPI } from '$lib/server/api';
import { error, fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

type Allocation = {
	id: string;
	chargeId: string;
	amountMinor: number;
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
	allocations: Allocation[];
};

type ReceiptAllocation = {
	chargeId: string;
	description: string;
	householdName: string;
	amountMinor: number;
};

type Receipt = {
	id: string;
	paymentId: string;
	receiptNumber: string;
	issuedAt: string;
	snapshot: {
		payerName: string;
		amountMinor: number;
		currency: string;
		method: string;
		reference?: string | null;
		allocations: ReceiptAllocation[];
		allocatedMinor: number;
		unallocatedMinor: number;
	};
};

type Contact = { id: string; displayName: string };

export const load: PageServerLoad = async ({ fetch, params }) => {
	const [paymentResult, receiptResult, contactResult] = await Promise.all([
		fetchAPI<Payment | null>(fetch, `/api/payments/${params.id}`, null),
		fetchAPI<Receipt | null>(fetch, `/api/payments/${params.id}/receipt`, null),
		fetchAPI<{ contacts: Contact[] }>(fetch, '/api/contacts', { contacts: [] })
	]);
	if (paymentResult.status === 404) error(404, 'Payment not found');
	if (!paymentResult.data) error(503, paymentResult.error ?? 'Payment data is unavailable');

	return {
		payment: paymentResult.data,
		receipt: receiptResult.status === 404 ? null : receiptResult.data,
		payerName: paymentResult.data.payerContactId
			? (contactResult.data.contacts.find(
					(contact) => contact.id === paymentResult.data?.payerContactId
				)?.displayName ?? 'Payer details are unavailable')
			: 'Payer not recorded',
		apiConnected:
			paymentResult.connected &&
			contactResult.connected &&
			(receiptResult.connected || receiptResult.status === 404)
	};
};

export const actions: Actions = {
	issueReceipt: async ({ fetch, params }) => {
		const result = await sendAPI<Receipt>(fetch, `/api/payments/${params.id}/receipt`, 'POST', {});
		if (result.error) return fail(result.status ?? 400, { error: result.error });
		throw redirect(303, `/payments/${params.id}`);
	}
};
