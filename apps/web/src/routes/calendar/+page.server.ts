import { dateParam, parseCalendarDate, parseCalendarView } from '$lib/dates/calendar';
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
		endAt?: string | null;
		locationType: string;
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
type Staff = { id: string; displayName: string };

const emptyDashboard: DashboardResponse = {
	dateStart: '',
	dateEnd: '',
	bookings: [],
	careTasks: [],
	openCharges: []
};

export const load: PageServerLoad = async ({ fetch, url }) => {
	const selectedDate = parseCalendarDate(url.searchParams.get('date'));
	const selectedDateParam = dateParam(selectedDate);
	const view = parseCalendarView(url.searchParams.get('view'));
	const [dashboard, households, staff] = await Promise.all([
		fetchAPI<DashboardResponse>(
			fetch,
			`/api/dashboard/today?date=${selectedDateParam}`,
			emptyDashboard
		),
		fetchAPI<{ households: Household[] }>(fetch, '/api/households', { households: [] }),
		fetchAPI<{ staff: Staff[] }>(fetch, '/api/staff', { staff: [] })
	]);

	return {
		apiConnected: dashboard.connected && households.connected && staff.connected,
		apiError: dashboard.error ?? households.error ?? staff.error,
		selectedDate: selectedDateParam,
		view,
		dashboard: dashboard.data,
		householdNames: Object.fromEntries(
			households.data.households.map((household) => [household.id, household.displayName])
		),
		staffNames: Object.fromEntries(
			staff.data.staff.map((member) => [member.id, member.displayName])
		)
	};
};
