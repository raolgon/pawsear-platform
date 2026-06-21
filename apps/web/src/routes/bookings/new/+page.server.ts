import { fail, redirect } from '@sveltejs/kit';
import { fetchAPI, sendAPI } from '$lib/server/api';
import type { Actions, PageServerLoad } from './$types';

type Household = {
	id: string;
	displayName: string;
	addressLine1?: string | null;
	neighborhood?: string | null;
	city?: string | null;
};

type Pet = {
	id: string;
	householdId: string;
	name: string;
	species: string;
	behaviorNotes?: string | null;
	feedingNotes?: string | null;
	medicalNotes?: string | null;
	active: boolean;
};

type Contact = {
	contactId: string;
	displayName: string;
	role: string;
	isPrimary: boolean;
};

type Staff = {
	id: string;
	displayName: string;
	role: string;
	active: boolean;
};

type CreatedBooking = {
	id?: string;
	message?: string;
};

const value = (formData: FormData, key: string) => {
	const raw = formData.get(key);
	return typeof raw === 'string' ? raw.trim() : '';
};

const optional = (formData: FormData, key: string) => {
	const raw = value(formData, key);
	return raw === '' ? undefined : raw;
};

const valuesFrom = (formData: FormData) => ({
	householdId: value(formData, 'householdId'),
	serviceType: value(formData, 'serviceType'),
	status: value(formData, 'status'),
	date: value(formData, 'date'),
	startTime: value(formData, 'startTime'),
	durationMinutes: value(formData, 'durationMinutes'),
	locationType: value(formData, 'locationType'),
	requestedByContactId: value(formData, 'requestedByContactId'),
	assignedStaffId: value(formData, 'assignedStaffId'),
	source: value(formData, 'source'),
	notes: value(formData, 'notes'),
	petIds: formData.getAll('petIds').filter((item): item is string => typeof item === 'string')
});

const localDateTimeToISO = (date: string, time: string) => {
	const parsed = new Date(`${date}T${time}:00`);
	return Number.isNaN(parsed.getTime()) ? '' : parsed.toISOString();
};

export const load: PageServerLoad = async ({ fetch, url }) => {
	const requestedPetId = url.searchParams.get('petId') ?? '';
	const requestedHouseholdId = url.searchParams.get('householdId') ?? '';
	const petResult = requestedPetId
		? await fetchAPI<Pet>(fetch, `/api/pets/${requestedPetId}`, null as unknown as Pet)
		: null;
	const selectedHouseholdId = requestedHouseholdId || petResult?.data?.householdId || '';

	const [households, pets, contacts, staff] = await Promise.all([
		fetchAPI<{ households: Household[] }>(fetch, '/api/households', { households: [] }),
		selectedHouseholdId
			? fetchAPI<{ pets: Pet[] }>(fetch, `/api/pets?householdId=${selectedHouseholdId}`, {
					pets: []
				})
			: fetchAPI<{ pets: Pet[] }>(fetch, '/api/pets', { pets: [] }),
		selectedHouseholdId
			? fetchAPI<{ contacts: Contact[] }>(
					fetch,
					`/api/households/${selectedHouseholdId}/contacts`,
					{ contacts: [] }
				)
			: Promise.resolve({ data: { contacts: [] }, connected: true }),
		fetchAPI<{ staff: Staff[] }>(fetch, '/api/staff', { staff: [] })
	]);

	const now = new Date();
	const roundedHour = new Date(now.getTime() + 60 * 60 * 1000);
	roundedHour.setMinutes(0, 0, 0);

	return {
		apiConnected: households.connected && pets.connected && staff.connected,
		households: households.data.households,
		pets: pets.data.pets,
		contacts: contacts.data.contacts,
		staff: staff.data.staff,
		selectedHouseholdId,
		selectedPetIds: requestedPetId ? [requestedPetId] : [],
		defaultDate: roundedHour.toISOString().slice(0, 10),
		defaultStartTime: roundedHour.toTimeString().slice(0, 5)
	};
};

export const actions: Actions = {
	default: async ({ request, fetch }) => {
		const formData = await request.formData();
		const values = valuesFrom(formData);
		const householdId = values.householdId;
		const serviceType = values.serviceType || 'walk';
		const status = values.status || 'requested';
		const locationType = values.locationType || 'household_home';
		const durationMinutes = Number.parseInt(values.durationMinutes || '45', 10);
		const startAt = localDateTimeToISO(values.date, values.startTime);

		if (!householdId) return fail(400, { values, error: 'Choose a household.' });
		if (!values.date || !values.startTime || !startAt)
			return fail(400, { values, error: 'Choose a valid date and time.' });
		if (!Number.isFinite(durationMinutes) || durationMinutes <= 0) {
			return fail(400, { values, error: 'Choose a valid duration.' });
		}

		const endAt = new Date(new Date(startAt).getTime() + durationMinutes * 60 * 1000).toISOString();
		const result = await sendAPI<CreatedBooking>(fetch, '/api/bookings', 'POST', {
			householdId,
			serviceType,
			status,
			startAt,
			endAt,
			locationType,
			requestedByContactId: optional(formData, 'requestedByContactId'),
			assignedStaffId: optional(formData, 'assignedStaffId'),
			source: values.source || 'manual',
			notes: optional(formData, 'notes'),
			petIds: values.petIds
		});

		if (!result.connected) return fail(503, { values, error: 'The API is not available.' });
		if (result.error || !result.data.id) {
			return fail(400, { values, error: result.error ?? 'Could not create booking.' });
		}

		throw redirect(303, `/households/${householdId}`);
	}
};
