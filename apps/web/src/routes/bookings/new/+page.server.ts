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

type DetectedRequest = {
	householdId?: string | null;
	contactId?: string | null;
	body: string;
	serviceType?: string | null;
	startAt?: string | null;
	endAt?: string | null;
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
	detectedRequestId: value(formData, 'detectedRequestId'),
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

const normalizedWords = (value: string) =>
	value
		.normalize('NFD')
		.replace(/\p{Diacritic}/gu, '')
		.toLocaleLowerCase('es')
		.split(/[^\p{Letter}\p{Number}]+/u)
		.filter(Boolean);

const messageMentionsPet = (body: string, petName: string) => {
	const message = normalizedWords(body);
	const name = normalizedWords(petName);
	return name.length > 0 && name.every((word) => message.includes(word));
};

export const load: PageServerLoad = async ({ fetch, url }) => {
	const detectedRequestId = url.searchParams.get('requestId') ?? '';
	const detectedRequestResult = detectedRequestId
		? await fetchAPI<DetectedRequest>(fetch, `/api/detected-requests/${detectedRequestId}`, {
				body: ''
			})
		: null;
	const detectedRequest = detectedRequestResult?.data;
	const requestedPetId = url.searchParams.get('petId') ?? '';
	const requestedHouseholdId = url.searchParams.get('householdId') ?? '';
	const petResult = requestedPetId
		? await fetchAPI<Pet>(fetch, `/api/pets/${requestedPetId}`, null as unknown as Pet)
		: null;
	const selectedHouseholdId =
		requestedHouseholdId || petResult?.data?.householdId || detectedRequest?.householdId || '';

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
	const detectedStartAt = detectedRequest?.startAt ?? '';
	const detectedEndAt = detectedRequest?.endAt ?? '';
	const detectedDuration =
		detectedStartAt && detectedEndAt
			? Math.round(
					(new Date(detectedEndAt).getTime() - new Date(detectedStartAt).getTime()) / 60000
				)
			: 45;
	const suggestedPetIds = detectedRequest?.body
		? pets.data.pets
				.filter((pet) => messageMentionsPet(detectedRequest.body, pet.name))
				.map((pet) => pet.id)
		: [];

	return {
		apiConnected: households.connected && pets.connected && staff.connected,
		households: households.data.households,
		pets: pets.data.pets,
		contacts: contacts.data.contacts,
		staff: staff.data.staff,
		detectedRequestId,
		selectedHouseholdId,
		selectedPetIds: requestedPetId ? [requestedPetId] : suggestedPetIds,
		defaultServiceType: detectedRequest?.serviceType ?? 'walk',
		defaultRequestedByContactId: detectedRequest?.contactId ?? '',
		defaultDate:
			(url.searchParams.get('date') ?? '') ||
			detectedStartAt.slice(0, 10) ||
			roundedHour.toISOString().slice(0, 10),
		defaultStartTime: detectedStartAt.slice(11, 16) || roundedHour.toTimeString().slice(0, 5),
		defaultDurationMinutes: [30, 45, 60, 90, 120, 1440].includes(detectedDuration)
			? String(detectedDuration)
			: '45'
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
		const bookingPath = values.detectedRequestId
			? `/api/detected-requests/${values.detectedRequestId}/bookings`
			: '/api/bookings';
		const result = await sendAPI<CreatedBooking>(fetch, bookingPath, 'POST', {
			householdId,
			serviceType,
			status,
			startAt,
			endAt,
			locationType,
			requestedByContactId: optional(formData, 'requestedByContactId'),
			assignedStaffId: optional(formData, 'assignedStaffId'),
			...(values.detectedRequestId ? {} : { source: values.source || 'manual' }),
			notes: optional(formData, 'notes'),
			petIds: values.petIds
		});

		if (!result.connected) return fail(503, { values, error: 'The API is not available.' });
		const bookingId =
			result.data.id ??
			('booking' in result.data && typeof result.data.booking === 'object' && result.data.booking
				? (result.data.booking as { id?: string }).id
				: undefined);
		if (result.error || !bookingId) {
			return fail(400, { values, error: result.error ?? 'Could not create booking.' });
		}

		if (values.detectedRequestId) {
			throw redirect(303, `/calendar?date=${values.date}`);
		}
		throw redirect(303, `/households/${householdId}`);
	}
};
