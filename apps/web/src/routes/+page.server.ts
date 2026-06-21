import { fetchAPI } from '$lib/server/api';
import type { PageServerLoad } from './$types';

type DashboardResponse = {
	dateStart: string;
	dateEnd: string;
	bookings: Array<{
		id: string;
		householdId: string;
		serviceType: string;
		status: string;
		startAt: string;
		assignedStaffId?: string | null;
		notes?: string | null;
	}>;
	careTasks: Array<{
		id: string;
		bookingId?: string | null;
		householdId: string;
		petId?: string | null;
		taskType: string;
		title: string;
		dueAt?: string | null;
		status: string;
	}>;
	openCharges: Array<{
		id: string;
		householdId: string;
		description: string;
		amountMinor: number;
		status: string;
		dueDate?: string | null;
	}>;
};

type Household = { id: string; displayName: string };
type Pet = { id: string; name: string };
type Staff = { id: string; displayName: string };

const emptyDashboard: DashboardResponse = {
	dateStart: '',
	dateEnd: '',
	bookings: [],
	careTasks: [],
	openCharges: []
};

export const load: PageServerLoad = async ({ fetch }) => {
	const [dashboard, households, pets, staff] = await Promise.all([
		fetchAPI<DashboardResponse>(fetch, '/api/dashboard/today', emptyDashboard),
		fetchAPI<{ households: Household[] }>(fetch, '/api/households', { households: [] }),
		fetchAPI<{ pets: Pet[] }>(fetch, '/api/pets', { pets: [] }),
		fetchAPI<{ staff: Staff[] }>(fetch, '/api/staff', { staff: [] })
	]);

	return {
		apiConnected: dashboard.connected && households.connected && pets.connected && staff.connected,
		apiError: dashboard.error ?? households.error ?? pets.error ?? staff.error,
		dashboard: dashboard.data,
		householdNames: Object.fromEntries(
			households.data.households.map((household) => [household.id, household.displayName])
		),
		petNames: Object.fromEntries(pets.data.pets.map((pet) => [pet.id, pet.name])),
		staffNames: Object.fromEntries(
			staff.data.staff.map((member) => [member.id, member.displayName])
		)
	};
};
