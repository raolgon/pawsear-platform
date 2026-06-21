import { fail, redirect } from '@sveltejs/kit';
import { fetchAPI, sendAPI } from '$lib/server/api';
import type { Actions, PageServerLoad } from './$types';

type Household = {
	id: string;
	displayName: string;
};

type CreatedPet = {
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
		const householdId = value(formData, 'householdId');
		const name = value(formData, 'name');
		const species = value(formData, 'species') || 'dog';
		const values = {
			householdId,
			name,
			species,
			breed: value(formData, 'breed'),
			size: value(formData, 'size'),
			sex: value(formData, 'sex'),
			birthdate: value(formData, 'birthdate'),
			colorMarkings: value(formData, 'colorMarkings'),
			behaviorNotes: value(formData, 'behaviorNotes'),
			medicalNotes: value(formData, 'medicalNotes'),
			feedingNotes: value(formData, 'feedingNotes'),
			vetNotes: value(formData, 'vetNotes')
		};

		if (!householdId) return fail(400, { values, error: 'Choose a household.' });
		if (!name) return fail(400, { values, error: 'Pet name is required.' });

		const result = await sendAPI<CreatedPet>(fetch, '/api/pets', 'POST', {
			householdId,
			name,
			species,
			breed: optional(formData, 'breed'),
			size: optional(formData, 'size'),
			sex: optional(formData, 'sex'),
			birthdate: optional(formData, 'birthdate'),
			colorMarkings: optional(formData, 'colorMarkings'),
			behaviorNotes: optional(formData, 'behaviorNotes'),
			medicalNotes: optional(formData, 'medicalNotes'),
			feedingNotes: optional(formData, 'feedingNotes'),
			vetNotes: optional(formData, 'vetNotes')
		});

		if (!result.connected) return fail(503, { values, error: 'The API is not available.' });
		if (result.error || !result.data.id) {
			return fail(400, { values, error: result.error ?? 'Could not create pet.' });
		}

		throw redirect(303, `/households/${householdId}`);
	}
};
