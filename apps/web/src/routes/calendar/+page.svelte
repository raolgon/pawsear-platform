<script lang="ts">
	import { resolve } from '$app/paths';
	import type { PageData } from './$types';
	import {
		AlertTriangle,
		CalendarDays,
		ChevronLeft,
		ChevronRight,
		Clock3,
		Footprints,
		Home,
		MoreVertical,
		Plus,
		Utensils
	} from 'lucide-svelte';
	import {
		calendarTitle,
		calendarViews,
		dateParam,
		monthGrid,
		nextDate,
		parseCalendarDate,
		previousDate,
		shortDateLabel,
		timeLabel,
		weekDays,
		type CalendarView
	} from '$lib/dates/calendar';

	type ScheduleItem = {
		id: string;
		time: string;
		endTime: string;
		title: string;
		meta: string;
		status: string;
		householdId: string;
		tone: 'booking' | 'task' | 'attention';
	};

	let { data }: { data: PageData } = $props();
	const selectedDate = $derived(parseCalendarDate(data.selectedDate));
	const title = $derived(calendarTitle(selectedDate, data.view));
	const selectedDay = $derived(data.selectedDate.slice(-2).replace(/^0/, ''));
	const selectedMonth = $derived(title.split(',').at(-1)?.trim() ?? title);
	const previousQuery = $derived(
		`?date=${dateParam(previousDate(selectedDate))}&view=${data.view}` as `?${string}`
	);
	const nextQuery = $derived(
		`?date=${dateParam(nextDate(selectedDate))}&view=${data.view}` as `?${string}`
	);
	const week = $derived(weekDays(selectedDate));
	const month = $derived(monthGrid(selectedDate));
	const calendarDays = $derived(data.view === 'month' ? month : week);
	const dayJobCount = $derived(data.dashboard.bookings.length + data.dashboard.careTasks.length);

	const viewQuery = (view: CalendarView): `?${string}` => `?date=${data.selectedDate}&view=${view}`;
	const dayQuery = (date: string): `?${string}` => `?date=${date}&view=${data.view}`;
	const createBookingQuery = (date = data.selectedDate): `?${string}` => `?date=${date}`;

	const items = $derived(
		[
			...data.dashboard.bookings.map(
				(booking): ScheduleItem => ({
					id: booking.id,
					time: timeLabel(booking.startAt),
					endTime: timeLabel(booking.endAt),
					title: booking.serviceType.replaceAll('_', ' '),
					meta: `${data.householdNames[booking.householdId] ?? 'Add household details'}${booking.assignedStaffId ? ` · ${data.staffNames[booking.assignedStaffId] ?? 'Assign a team member'}` : ' · Assign a team member'}`,
					status: booking.status,
					householdId: booking.householdId,
					tone: booking.status === 'requested' ? 'attention' : 'booking'
				})
			),
			...data.dashboard.careTasks.map(
				(task): ScheduleItem => ({
					id: task.id,
					time: timeLabel(task.dueAt),
					endTime: '',
					title: task.title,
					meta: `${task.taskType.replaceAll('_', ' ')} · ${data.householdNames[task.householdId] ?? 'Add household details'}`,
					status: task.status,
					householdId: task.householdId,
					tone: task.status === 'pending' ? 'task' : 'booking'
				})
			)
		].sort((a, b) => a.time.localeCompare(b.time))
	);

	const itemClasses = (tone: ScheduleItem['tone'], highlighted: boolean) => {
		if (highlighted)
			return 'border-prussian-blue-700 bg-prussian-blue-700 text-white shadow-lg shadow-prussian-blue-900/15';
		if (tone === 'attention')
			return 'border-sandy-brown-300 bg-sandy-brown-50 text-sandy-brown-900';
		if (tone === 'task') return 'border-light-gold-300 bg-white text-prussian-blue-950';
		return 'border-prussian-blue-100 bg-prussian-blue-50 text-prussian-blue-950';
	};
</script>

<svelte:head>
	<title>Experimental Calendar · Pawsear</title>
</svelte:head>

<main class="min-h-screen bg-[#faf9f1] text-prussian-blue-950">
	<div
		class="mx-auto grid w-full max-w-6xl grid-cols-1 lg:grid-cols-[minmax(0,44rem)_20rem] lg:gap-8 lg:px-8"
	>
		<section class="min-w-0">
			<header class="px-4 pt-5 pb-3 lg:px-0 lg:pt-7">
				<div class="flex items-center justify-between gap-3">
					<div class="min-w-0">
						<p class="text-xs font-bold tracking-wide text-light-gold-800 uppercase">
							Experimental calendar
						</p>
						<div class="mt-1 flex items-end gap-2">
							<span class="text-4xl leading-none font-black">{selectedDay}</span>
							<span class="pb-1 text-sm font-black text-light-gold-800 uppercase"
								>{selectedMonth}</span
							>
						</div>
					</div>
					<button
						class="flex h-10 w-10 items-center justify-center rounded-lg border border-light-gold-200 bg-white text-prussian-blue-800 shadow-sm shadow-black/5"
						aria-label="Calendar options"
					>
						<MoreVertical size={18} />
					</button>
				</div>
			</header>

			<div class="space-y-4 px-4 pb-24 lg:px-0 lg:pb-10">
				{#if !data.apiConnected}
					<section
						class="rounded-lg border border-sandy-brown-200 bg-sandy-brown-50 p-3 text-sm font-bold text-sandy-brown-800"
					>
						API offline. Calendar is showing an empty experimental state.
					</section>
				{/if}

				<section class="rounded-lg border border-light-gold-200 bg-white shadow-sm shadow-black/5">
					<div
						class="flex items-center justify-between gap-2 border-b border-light-gold-100 px-3 py-3"
					>
						<a
							href={resolve(`/calendar${previousQuery}`)}
							class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg border border-light-gold-200 text-prussian-blue-800"
							aria-label="Previous day"
						>
							<ChevronLeft size={19} />
						</a>

						<form method="GET" class="grid min-w-0 flex-1 grid-cols-[minmax(0,1fr)_auto] gap-2">
							<input type="hidden" name="view" value={data.view} />
							<input
								name="date"
								type="date"
								value={data.selectedDate}
								class="min-w-0 rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-2 text-center text-sm font-black outline-none focus:border-prussian-blue-700"
							/>
							<button
								class="rounded-lg bg-prussian-blue-700 px-3 py-2 text-sm font-black text-white"
								>Go</button
							>
						</form>

						<a
							href={resolve(`/calendar${nextQuery}`)}
							class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg border border-light-gold-200 text-prussian-blue-800"
							aria-label="Next day"
						>
							<ChevronRight size={19} />
						</a>
					</div>

					<div class="px-3 py-3">
						<div class="mb-3 grid grid-cols-3 rounded-lg bg-[#f1ead7] p-1 text-sm font-bold">
							{#each calendarViews as view (view)}
								<a
									href={resolve(`/calendar${viewQuery(view)}`)}
									class={[
										'rounded-md px-3 py-2 text-center capitalize',
										data.view === view
											? 'bg-white text-prussian-blue-950 shadow-sm'
											: 'text-light-gold-800'
									]}
								>
									{view}
								</a>
							{/each}
						</div>

						{#if data.view === 'month'}
							<div
								class="mb-2 grid grid-cols-7 gap-1 text-center text-[0.68rem] font-black text-light-gold-800 uppercase"
							>
								{#each ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'] as label (label)}
									<span>{label}</span>
								{/each}
							</div>
							<div class="grid grid-cols-7 gap-y-3">
								{#each calendarDays as day (day.dateParam)}
									<a
										href={resolve(`/calendar${dayQuery(day.dateParam)}`)}
										class={[
											'flex min-w-0 flex-col items-center justify-center gap-1 text-sm font-black',
											day.isSelected
												? 'text-prussian-blue-950'
												: day.isCurrentMonth
													? 'text-prussian-blue-950'
													: 'text-light-gold-400'
										]}
										aria-label={day.dateParam}
									>
										<span
											class={[
												'flex h-9 w-9 items-center justify-center rounded-full',
												day.isSelected ? 'bg-prussian-blue-700 text-white shadow-sm' : ''
											]}
										>
											{day.dayLabel}
										</span>
										<span
											class={[
												'h-1 w-1 rounded-full',
												day.isSelected && dayJobCount > 0
													? 'bg-prussian-blue-700'
													: 'bg-transparent'
											]}
										></span>
									</a>
								{/each}
							</div>
						{:else}
							<div class="grid grid-cols-7 gap-y-3">
								{#each calendarDays as day (day.dateParam)}
									<a
										href={resolve(`/calendar${dayQuery(day.dateParam)}`)}
										class={[
											'min-w-0 text-center',
											day.isSelected ? 'text-prussian-blue-950' : 'text-prussian-blue-950'
										]}
										aria-label={day.dateParam}
									>
										<span
											class="block truncate text-[0.68rem] font-black text-light-gold-800 uppercase"
											>{day.weekdayLabel.slice(0, 3)}</span
										>
										<span
											class={[
												'mx-auto mt-1 flex h-9 w-9 items-center justify-center rounded-full text-base font-black',
												day.isSelected ? 'bg-prussian-blue-700 text-white shadow-sm' : ''
											]}>{day.dayLabel}</span
										>
										<span
											class={[
												'mx-auto mt-1 block h-1 w-1 rounded-full',
												day.isSelected && dayJobCount > 0
													? 'bg-prussian-blue-700'
													: 'bg-transparent'
											]}
										></span>
									</a>
								{/each}
							</div>
						{/if}
					</div>
				</section>

				<section class="space-y-3">
					<div
						class="flex items-center justify-between rounded-lg border border-light-gold-200 bg-white px-4 py-3 shadow-sm shadow-black/5"
					>
						<div>
							<p class="text-sm font-black">Today</p>
							<p class="text-xs font-bold text-light-gold-800 uppercase">
								{shortDateLabel(selectedDate)}
							</p>
						</div>
						<span class="text-sm font-black text-prussian-blue-800">{items.length} events</span>
					</div>

					{#if items.length === 0}
						<div
							class="rounded-lg border border-dashed border-light-gold-300 bg-white p-5 text-center shadow-sm shadow-black/5"
						>
							<div
								class="mx-auto mb-3 flex h-11 w-11 items-center justify-center rounded-lg bg-vanilla-custard-100 text-prussian-blue-800"
							>
								<Clock3 size={21} />
							</div>
							<h3 class="text-base font-black">No scheduled work</h3>
							<p class="mx-auto mt-1 max-w-sm text-sm leading-6 font-semibold text-light-gold-800">
								Create a booking for this day or use the dashboard to review pending work.
							</p>
							<a
								href={resolve(`/bookings/new${createBookingQuery()}`)}
								class="mt-4 inline-flex items-center justify-center gap-2 rounded-lg bg-prussian-blue-700 px-4 py-3 text-sm font-black text-white"
							>
								<Plus size={18} />
								Create booking
							</a>
						</div>
					{:else}
						<div class="space-y-2">
							{#each items as item, index (item.id)}
								{@const highlighted = item.tone === 'attention' || index === 0}
								<a
									href={resolve('/households/[slug]', { slug: item.householdId })}
									class="relative block"
								>
									<article
										class={[
											'rounded-lg border p-4 shadow-sm shadow-black/5',
											itemClasses(item.tone, highlighted)
										]}
									>
										<div class="grid grid-cols-[4.25rem_minmax(0,1fr)] gap-3">
											<div
												class={[
													'border-r pr-3',
													highlighted ? 'border-white/25' : 'border-light-gold-200'
												]}
											>
												<p class="text-sm font-black">{item.time}</p>
												{#if item.endTime && item.endTime !== 'Any time'}
													<p
														class={[
															'mt-1 text-xs font-bold',
															highlighted ? 'text-white/70' : 'text-light-gold-800'
														]}
													>
														{item.endTime}
													</p>
												{/if}
											</div>
											<div class="flex min-w-0 items-start gap-3">
												<div
													class={[
														'flex h-10 w-10 shrink-0 items-center justify-center rounded-lg',
														highlighted
															? 'bg-white/15 text-white'
															: 'bg-white text-prussian-blue-800'
													]}
												>
													{#if item.tone === 'task'}
														<Utensils size={19} />
													{:else if item.tone === 'attention'}
														<AlertTriangle size={19} />
													{:else}
														<Footprints size={19} />
													{/if}
												</div>
												<div class="min-w-0 flex-1">
													<div class="flex items-start justify-between gap-3">
														<div class="min-w-0">
															<h3 class="truncate text-base font-black capitalize">{item.title}</h3>
															<p
																class={[
																	'mt-1 truncate text-xs font-bold',
																	highlighted ? 'text-white/80' : 'text-light-gold-800'
																]}
															>
																{item.meta}
															</p>
														</div>
														<span
															class={[
																'shrink-0 rounded-full px-2 py-1 text-xs font-bold ring-1',
																highlighted
																	? 'bg-white/15 text-white ring-white/20'
																	: 'bg-white text-prussian-blue-800 ring-light-gold-200'
															]}>{item.status}</span
														>
													</div>
												</div>
											</div>
										</div>
									</article>
								</a>
							{/each}
						</div>
					{/if}
				</section>
			</div>
		</section>

		<aside class="hidden border-l border-light-gold-200 bg-[#f6f3e8]/80 px-5 py-7 lg:block">
			<div class="sticky top-24 space-y-5">
				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
				>
					<div class="mb-3 flex items-center gap-2">
						<CalendarDays size={18} class="text-prussian-blue-700" />
						<h2 class="text-base font-black">Date state</h2>
					</div>
					<div class="space-y-3 text-sm">
						<div>
							<p class="text-xs font-bold text-light-gold-800 uppercase">Selected</p>
							<p class="font-black">{data.selectedDate}</p>
						</div>
						<div>
							<p class="text-xs font-bold text-light-gold-800 uppercase">View</p>
							<p class="font-black capitalize">{data.view}</p>
						</div>
						<div>
							<p class="text-xs font-bold text-light-gold-800 uppercase">Loaded range</p>
							<p class="font-black">
								{data.dashboard.dateStart
									? `${timeLabel(data.dashboard.dateStart)}-${timeLabel(data.dashboard.dateEnd)}`
									: 'Start the local API to load dates'}
							</p>
						</div>
					</div>
				</section>

				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 text-sm shadow-sm shadow-black/5"
				>
					<div class="mb-2 flex items-center gap-2">
						<Home size={18} class="text-prussian-blue-700" />
						<h2 class="text-base font-black">Experiment notes</h2>
					</div>
					<p class="leading-6 font-semibold text-light-gold-800">
						The layout is inspired by compact mobile calendars, but stays in Pawsear's light
						operational style.
					</p>
				</section>
			</div>
		</aside>
	</div>

	<a
		href={resolve(`/bookings/new${createBookingQuery()}`)}
		class="fixed right-4 bottom-5 z-20 flex h-14 w-14 items-center justify-center rounded-full bg-prussian-blue-700 text-white shadow-xl shadow-prussian-blue-950/20 lg:hidden"
		aria-label="Create booking"
	>
		<Plus size={25} />
	</a>
</main>
