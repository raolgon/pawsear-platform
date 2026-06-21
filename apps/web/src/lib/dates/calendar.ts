import {
	addDays,
	eachDayOfInterval,
	endOfMonth,
	format,
	isSameDay,
	isSameMonth,
	isValid,
	parseISO,
	startOfMonth,
	startOfWeek,
	subDays
} from 'date-fns';

export type CalendarView = 'day' | 'week' | 'month';

export type CalendarDay = {
	date: Date;
	dateParam: string;
	dayLabel: string;
	weekdayLabel: string;
	isSelected: boolean;
	isCurrentMonth: boolean;
};

export const calendarViews: CalendarView[] = ['day', 'week', 'month'];

export function parseCalendarView(value: string | null): CalendarView {
	return value === 'week' || value === 'month' ? value : 'day';
}

export function dateParam(date: Date): string {
	return format(date, 'yyyy-MM-dd');
}

export function parseCalendarDate(value: string | null, fallback = new Date()): Date {
	if (!value) return fallback;
	const parsed = parseISO(value);
	return isValid(parsed) ? parsed : fallback;
}

export function nextDate(date: Date): Date {
	return addDays(date, 1);
}

export function previousDate(date: Date): Date {
	return subDays(date, 1);
}

export function calendarTitle(date: Date, view: CalendarView): string {
	if (view === 'month') return format(date, 'MMMM yyyy');
	if (view === 'week') {
		const start = startOfWeek(date, { weekStartsOn: 1 });
		const end = addDays(start, 6);
		return `${format(start, 'MMM d')} - ${format(end, 'MMM d, yyyy')}`;
	}
	return format(date, 'EEEE, MMMM d');
}

export function shortDateLabel(date: Date): string {
	return format(date, 'MMM d');
}

export function timeLabel(value?: string | null): string {
	if (!value) return 'Any time';
	const parsed = parseISO(value);
	return isValid(parsed) ? format(parsed, 'HH:mm') : value.slice(11, 16) || value;
}

export function weekDays(date: Date): CalendarDay[] {
	const start = startOfWeek(date, { weekStartsOn: 1 });
	return eachDayOfInterval({ start, end: addDays(start, 6) }).map((day) =>
		dayModel(day, date, date)
	);
}

export function monthGrid(date: Date): CalendarDay[] {
	const monthStart = startOfMonth(date);
	const monthEnd = endOfMonth(date);
	const gridStart = startOfWeek(monthStart, { weekStartsOn: 1 });
	const gridEnd = addDays(startOfWeek(monthEnd, { weekStartsOn: 1 }), 6);
	return eachDayOfInterval({ start: gridStart, end: gridEnd }).map((day) =>
		dayModel(day, date, monthStart)
	);
}

function dayModel(day: Date, selectedDate: Date, visibleMonth: Date): CalendarDay {
	return {
		date: day,
		dateParam: dateParam(day),
		dayLabel: format(day, 'd'),
		weekdayLabel: format(day, 'EEE'),
		isSelected: isSameDay(day, selectedDate),
		isCurrentMonth: isSameMonth(day, visibleMonth)
	};
}
