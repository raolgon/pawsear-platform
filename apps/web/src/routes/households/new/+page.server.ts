import { fail, redirect } from '@sveltejs/kit';
import { sendAPI } from '$lib/server/api';
import type { Actions } from './$types';

type CreatedHousehold = {
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

export const actions: Actions = {
	default: async ({ request, fetch }) => {
		const formData = await request.formData();
		const displayName = value(formData, 'displayName');

		const values = {
			displayName,
			addressLine1: value(formData, 'addressLine1'),
			addressLine2: value(formData, 'addressLine2'),
			neighborhood: value(formData, 'neighborhood'),
			city: value(formData, 'city'),
			timezone: value(formData, 'timezone'),
			notes: value(formData, 'notes')
		};

		if (!displayName) {
			return fail(400, {
				values,
				error: 'Household name is required.'
			});
		}

		const result = await sendAPI<CreatedHousehold>(fetch, '/api/households', 'POST', {
			displayName,
			addressLine1: optional(formData, 'addressLine1'),
			addressLine2: optional(formData, 'addressLine2'),
			neighborhood: optional(formData, 'neighborhood'),
			city: optional(formData, 'city'),
			timezone: optional(formData, 'timezone') ?? 'America/Mexico_City',
			notes: optional(formData, 'notes')
		});

		if (!result.connected) {
			return fail(503, {
				values,
				error: 'The API is not available. Start the backend and try again.'
			});
		}

		if (result.error || !result.data.id) {
			return fail(400, {
				values,
				error: result.error ?? 'Could not create household.'
			});
		}

		throw redirect(303, `/households/${result.data.id}`);
	}
};
