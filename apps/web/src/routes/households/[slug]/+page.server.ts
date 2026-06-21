import { fetchAPI, sendAPI } from '$lib/server/api';
import { error, fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

type Household = {
	id: string;
	displayName: string;
	neighborhood?: string | null;
	city?: string | null;
	notes?: string | null;
	active: boolean;
};

type Pet = {
	id: string;
	name: string;
	species: string;
	size?: string | null;
	behaviorNotes?: string | null;
	medicalNotes?: string | null;
	feedingNotes?: string | null;
	active: boolean;
};

type Contact = {
	id: string;
	displayName: string;
	role: string;
	isPrimary: boolean;
	phone?: string | null;
};

type Booking = {
	id: string;
	householdId: string;
	serviceType: string;
	status: string;
	startAt: string;
};

type CareTask = {
	id: string;
	householdId: string;
	title: string;
	taskType: string;
	status: string;
	dueAt?: string | null;
};

type Charge = {
	id: string;
	bookingId?: string | null;
	description: string;
	amountMinor: number;
	allocatedMinor: number;
	outstandingMinor: number;
	currency: string;
	status: string;
};

const money = (amountMinor: number, currency = 'MXN') =>
	new Intl.NumberFormat('en-US', { style: 'currency', currency, maximumFractionDigits: 0 }).format(
		amountMinor / 100
	);

export const load: PageServerLoad = async ({ fetch, params }) => {
	const [householdResult, petsResult, contactsResult, bookingsResult, tasksResult, chargesResult] =
		await Promise.all([
			fetchAPI<Household | null>(fetch, `/api/households/${params.slug}`, null),
			fetchAPI<{ pets: Pet[] }>(fetch, `/api/pets?householdId=${params.slug}`, { pets: [] }),
			fetchAPI<{ contacts: Contact[] }>(fetch, `/api/households/${params.slug}/contacts`, {
				contacts: []
			}),
			fetchAPI<{ bookings: Booking[] }>(fetch, `/api/bookings?householdId=${params.slug}`, {
				bookings: []
			}),
			fetchAPI<{ careTasks: CareTask[] }>(fetch, `/api/care-tasks?householdId=${params.slug}`, {
				careTasks: []
			}),
			fetchAPI<{ charges: Charge[] }>(fetch, `/api/charges?householdId=${params.slug}`, {
				charges: []
			})
		]);

	if (householdResult.status === 404) {
		error(404, 'Household not found');
	}
	if (
		!householdResult.connected ||
		!householdResult.data ||
		!petsResult.connected ||
		!contactsResult.connected ||
		!bookingsResult.connected ||
		!tasksResult.connected ||
		!chargesResult.connected
	) {
		error(503, 'The local Pawsear API is unavailable');
	}

	const household = householdResult.data;
	const pets = petsResult.data.pets;
	const contacts = contactsResult.data.contacts;
	const bookings = bookingsResult.data.bookings;
	const careTasks = tasksResult.data.careTasks;
	const charges = chargesResult.data.charges;
	const openBalance = charges
		.filter((charge) => charge.status === 'unpaid' || charge.status === 'partially_paid')
		.reduce((total, charge) => total + charge.outstandingMinor, 0);
	const nextBooking = bookings[0];

	return {
		apiConnected: true,
		bookings,
		careTasks,
		charges,
		contacts,
		household: {
			id: household.id,
			name: household.displayName,
			area: [household.neighborhood, household.city].filter(Boolean).join(' · ') || 'No area yet',
			status: household.active ? 'Active' : 'Inactive',
			alerts: `${bookings.length} bookings · ${charges.length} charges`,
			balance: openBalance > 0 ? `${money(openBalance)} unpaid` : 'No unpaid balance',
			nextWork: nextBooking
				? {
						time: `${nextBooking.startAt.slice(11, 16)} · ${nextBooking.serviceType.replaceAll('_', ' ')}`,
						pets: pets.map((pet) => pet.name).join(' + ') || household.displayName,
						status: nextBooking.status
					}
				: {
						time: 'No scheduled work',
						pets: pets.map((pet) => pet.name).join(' + ') || 'No pets yet',
						status: 'Open'
					},
			pets: pets.map((pet) => ({
				id: pet.id,
				name: pet.name,
				kind: `${pet.species}${pet.size ? ` · ${pet.size}` : ''}${pet.active ? ' · active' : ''}`,
				next: nextBooking
					? `${nextBooking.serviceType.replaceAll('_', ' ')} today`
					: 'No upcoming work',
				care: pet.feedingNotes ?? 'No feeding note yet',
				note: pet.behaviorNotes ?? pet.medicalNotes ?? 'No special notes yet'
			})),
			contacts: contacts.map((contact) => ({
				name: contact.displayName,
				role: `${contact.role}${contact.isPrimary ? ' · primary' : ''}${contact.phone ? ' · phone' : ''}`
			})),
			notes: household.notes ?? 'No household notes yet.',
			careFood: pets
				.map((pet) => `${pet.name}: ${pet.feedingNotes ?? 'No feeding note'}`)
				.join(' · '),
			careMedicine: pets
				.map((pet) => `${pet.name}: ${pet.medicalNotes ?? 'No medicine note'}`)
				.join(' · '),
			money:
				charges.length > 0
					? charges
							.map(
								(charge) =>
									`${charge.description}: ${money(charge.amountMinor, charge.currency)} ${charge.status}`
							)
							.join(' · ')
					: 'No charges yet.',
			history:
				bookings.length > 0
					? `${bookings.length} bookings loaded from API.`
					: 'No booking history yet.',
			messages: 'Message import is not connected yet.'
		}
	};
};

const positiveAmountMinor = (value: FormDataEntryValue | null): number | null => {
	if (typeof value !== 'string' || !/^\d+(?:\.\d{1,2})?$/.test(value.trim())) return null;
	const amountMinor = Math.round(Number(value) * 100);
	return Number.isSafeInteger(amountMinor) && amountMinor > 0 ? amountMinor : null;
};

export const actions: Actions = {
	bookingStatus: async ({ fetch, params, request }) => {
		const formData = await request.formData();
		const bookingId = String(formData.get('bookingId') ?? '');
		const status = String(formData.get('status') ?? '');
		if (!bookingId || !['confirmed', 'in_progress', 'completed', 'cancelled'].includes(status)) {
			return fail(400, { error: 'Choose a valid booking status.' });
		}
		const booking = await fetchAPI<Booking | null>(fetch, `/api/bookings/${bookingId}`, null);
		if (!booking.connected || !booking.data || booking.data.householdId !== params.slug) {
			return fail(404, { error: 'Booking not found for this household.' });
		}
		const timestamp = new Date().toISOString();
		const result = await sendAPI<Booking>(fetch, `/api/bookings/${bookingId}`, 'PATCH', {
			status,
			completedAt: status === 'completed' ? timestamp : undefined,
			cancelledAt: status === 'cancelled' ? timestamp : undefined
		});
		if (result.error) return fail(400, { error: result.error });
		throw redirect(303, `/households/${params.slug}`);
	},
	careTaskStatus: async ({ fetch, params, request }) => {
		const formData = await request.formData();
		const taskId = String(formData.get('taskId') ?? '');
		const status = String(formData.get('status') ?? '');
		const skippedReason = String(formData.get('skippedReason') ?? '').trim();
		if (!taskId || !['completed', 'skipped'].includes(status)) {
			return fail(400, { error: 'Choose a valid care task status.' });
		}
		if (status === 'skipped' && !skippedReason) {
			return fail(400, { error: 'A skipped task needs a reason.' });
		}
		const task = await fetchAPI<CareTask | null>(fetch, `/api/care-tasks/${taskId}`, null);
		if (!task.connected || !task.data || task.data.householdId !== params.slug) {
			return fail(404, { error: 'Care task not found for this household.' });
		}
		const result = await sendAPI<CareTask>(fetch, `/api/care-tasks/${taskId}`, 'PATCH', {
			status,
			completedAt: status === 'completed' ? new Date().toISOString() : undefined,
			skippedReason: status === 'skipped' ? skippedReason : undefined
		});
		if (result.error) return fail(400, { error: result.error });
		throw redirect(303, `/households/${params.slug}`);
	},
	createCharge: async ({ fetch, params, request }) => {
		const formData = await request.formData();
		const description = String(formData.get('description') ?? '').trim();
		const bookingId = String(formData.get('bookingId') ?? '').trim();
		const amountMinor = positiveAmountMinor(formData.get('amount'));
		if (!description || amountMinor === null) {
			return fail(400, { error: 'Charge description and a positive amount are required.' });
		}
		if (bookingId) {
			const booking = await fetchAPI<Booking | null>(fetch, `/api/bookings/${bookingId}`, null);
			if (!booking.connected || !booking.data || booking.data.householdId !== params.slug) {
				return fail(400, { error: 'The selected booking does not belong to this household.' });
			}
		}
		const result = await sendAPI<Charge>(fetch, '/api/charges', 'POST', {
			householdId: params.slug,
			bookingId: bookingId || undefined,
			description,
			amountMinor,
			currency: 'MXN',
			status: 'unpaid',
			dueDate: String(formData.get('dueDate') ?? '') || undefined
		});
		if (result.error) return fail(400, { error: result.error });
		throw redirect(303, `/households/${params.slug}`);
	}
};
