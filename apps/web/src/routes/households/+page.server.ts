import { fetchAPI } from '$lib/server/api';
import type { PageServerLoad } from './$types';

type Household = {
	id: string;
	displayName: string;
	neighborhood?: string | null;
	city?: string | null;
	active: boolean;
};

export const load: PageServerLoad = async ({ fetch }) => {
	const result = await fetchAPI<{ households: Household[] }>(fetch, '/api/households', {
		households: []
	});

	return {
		apiConnected: result.connected,
		apiError: result.error,
		households: result.data.households
	};
};
