import { fail, redirect } from '@sveltejs/kit';
import { fetchAPI, sendAPI } from '$lib/server/api';
import type { Actions, PageServerLoad } from './$types';

type Household = {
	id: string;
	displayName: string;
};

type CreatedContact = {
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

export const load: PageServerLoad = async ({ fetch, url }) => {
	const households = await fetchAPI<{ households: Household[] }>(fetch, '/api/households', {
		households: []
	});

	return {
		apiConnected: households.connected,
		households: households.data.households,
		selectedHouseholdId: url.searchParams.get('householdId') ?? ''
	};
};

export const actions: Actions = {
	default: async ({ request, fetch }) => {
		const formData = await request.formData();
		const displayName = value(formData, 'displayName');
		const householdId = value(formData, 'householdId');
		const role = value(formData, 'role') || 'owner';
		const values = {
			displayName,
			householdId,
			role,
			phone: value(formData, 'phone'),
			whatsappId: value(formData, 'whatsappId'),
			telegramId: value(formData, 'telegramId'),
			email: value(formData, 'email'),
			notes: value(formData, 'notes'),
			relationshipNotes: value(formData, 'relationshipNotes')
		};

		if (!displayName) return fail(400, { values, error: 'Contact name is required.' });
		if (!householdId) return fail(400, { values, error: 'Choose a household.' });

		const contact = await sendAPI<CreatedContact>(fetch, '/api/contacts', 'POST', {
			displayName,
			phone: optional(formData, 'phone'),
			whatsappId: optional(formData, 'whatsappId'),
			telegramId: optional(formData, 'telegramId'),
			email: optional(formData, 'email'),
			notes: optional(formData, 'notes')
		});

		if (!contact.connected) return fail(503, { values, error: 'The API is not available.' });
		if (contact.error || !contact.data.id) {
			return fail(400, { values, error: contact.error ?? 'Could not create contact.' });
		}

		const link = await sendAPI(fetch, `/api/households/${householdId}/contacts`, 'POST', {
			contactId: contact.data.id,
			role,
			isPrimary: formData.get('isPrimary') === 'on',
			notes: optional(formData, 'relationshipNotes')
		});

		if (link.error) return fail(400, { values, error: link.error });

		throw redirect(303, `/households/${householdId}`);
	}
};
