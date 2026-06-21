<script lang="ts">
	import { resolve } from '$app/paths';
	import type { PageData } from './$types';
	import {
		AlertCircle,
		CalendarDays,
		ChevronRight,
		CircleCheck,
		Clock3,
		Footprints,
		Home,
		Pill,
		Plus,
		SlidersHorizontal,
		Utensils
	} from 'lucide-svelte';

	type Tone = 'attention' | 'upcoming' | 'confirmed' | 'active' | 'complete';
	type Kind = 'medicine' | 'food' | 'walk' | 'boarding' | 'visit';
	type ScheduleItem = {
		id: string;
		householdId: string;
		time: string;
		title: string;
		detail: string;
		household: string;
		status: string;
		tone: Tone;
		kind: Kind;
	};

	let { data }: { data: PageData } = $props();

	const timeLabel = (value?: string | null) => value?.slice(11, 16) || 'Any time';
	const titleCase = (value: string) =>
		value.replaceAll('_', ' ').replace(/\b\w/g, (letter) => letter.toUpperCase());
	const selectedDate = $derived(
		data.dashboard.dateStart?.slice(0, 10) || 'Load the schedule to choose a date'
	);

	const schedule = $derived(
		[
			...data.dashboard.bookings.map(
				(booking): ScheduleItem => ({
					id: booking.id,
					householdId: booking.householdId,
					time: timeLabel(booking.startAt),
					title: titleCase(booking.serviceType),
					detail: booking.assignedStaffId
						? (data.staffNames[booking.assignedStaffId] ?? 'Add team member details')
						: 'Assign a team member',
					household: data.householdNames[booking.householdId] ?? 'Add household details',
					status: titleCase(booking.status),
					tone:
						booking.status === 'requested'
							? 'attention'
							: booking.status === 'completed'
								? 'complete'
								: booking.status === 'in_progress'
									? 'active'
									: 'confirmed',
					kind:
						booking.serviceType === 'boarding'
							? 'boarding'
							: booking.serviceType === 'walk'
								? 'walk'
								: 'visit'
				})
			),
			...data.dashboard.careTasks.map(
				(task): ScheduleItem => ({
					id: task.id,
					householdId: task.householdId,
					time: timeLabel(task.dueAt),
					title: task.petId
						? `${data.petNames[task.petId] ?? 'Add pet details'} · ${task.title}`
						: task.title,
					detail: titleCase(task.taskType),
					household: data.householdNames[task.householdId] ?? 'Add household details',
					status: titleCase(task.status),
					tone: task.status === 'pending' ? 'attention' : 'complete',
					kind:
						task.taskType === 'medicine' ? 'medicine' : task.taskType === 'food' ? 'food' : 'visit'
				})
			)
		].sort((a, b) => a.time.localeCompare(b.time))
	);

	const attention = $derived(schedule.filter((item) => item.tone === 'attention'));
	const boardingCount = $derived(schedule.filter((item) => item.kind === 'boarding').length);
	const completedCount = $derived(schedule.filter((item) => item.tone === 'complete').length);

	const iconFor = (kind: Kind) => {
		if (kind === 'medicine') return Pill;
		if (kind === 'food') return Utensils;
		if (kind === 'walk') return Footprints;
		if (kind === 'boarding') return Home;
		return Clock3;
	};

	const badgeClass = (tone: Tone) => {
		if (tone === 'attention') return 'bg-[#fff0ed] text-[#d84832]';
		if (tone === 'active') return 'bg-[#edf4ff] text-[#1c5c98]';
		if (tone === 'complete') return 'bg-[#eaf5e6] text-[#397546]';
		return 'bg-[#fbf3df] text-[#a76f12]';
	};

	const iconClass = (tone: Tone) => {
		if (tone === 'attention') return 'bg-[#fff0ed] text-[#d84832]';
		if (tone === 'complete') return 'bg-[#eaf5e6] text-[#397546]';
		return 'bg-[#fbf3df] text-[#b37b1c]';
	};
</script>

<svelte:head><title>Today · Pawsear</title></svelte:head>

<main class="min-h-screen bg-[#f7f6f2] text-[#071b3b]">
	<header class="border-b border-[#e9e3d6] bg-[#fffdf8]">
		<div
			class="mx-auto flex max-w-[1280px] items-center justify-between gap-4 px-4 py-5 sm:px-7 lg:px-10 lg:py-6"
		>
			<div>
				<p class="text-sm font-bold text-[#7a8797]">{selectedDate}</p>
				<h1 class="mt-0.5 text-3xl font-black tracking-[-0.04em] lg:text-4xl">Today</h1>
			</div>
			<div class="flex items-center gap-2">
				<button
					class="flex h-11 w-11 items-center justify-center rounded-xl border border-[#e2ddd2] bg-white text-[#526174]"
					aria-label="Filters"
				>
					<SlidersHorizontal size={19} />
				</button>
				<a
					href={resolve('/bookings/new')}
					class="flex h-11 items-center gap-2 rounded-xl bg-[#c58b20] px-4 text-sm font-black text-white shadow-sm shadow-[#c58b20]/20"
				>
					<Plus size={18} />
					<span class="hidden sm:inline">New booking</span>
				</a>
			</div>
		</div>
	</header>

	<div
		class="mx-auto grid max-w-[1280px] gap-7 px-4 py-6 sm:px-7 lg:grid-cols-[minmax(0,1fr)_310px] lg:px-10 lg:py-8"
	>
		<section class="min-w-0">
			{#if !data.apiConnected}
				<div
					class="mb-5 rounded-2xl border border-[#f3c7bf] bg-[#fff4f1] p-4 text-sm font-bold text-[#a63e2f]"
				>
					API offline · operational data is unavailable
				</div>
			{/if}

			<div class="mb-4 flex items-center justify-between">
				<div>
					<p class="text-xs font-black tracking-[0.14em] text-[#8b96a3] uppercase">Schedule</p>
					<h2 class="mt-1 text-xl font-black">Upcoming</h2>
				</div>
				<a
					href={resolve('/calendar')}
					class="flex items-center gap-2 rounded-xl border border-[#e2ddd2] bg-white px-3 py-2 text-sm font-extrabold text-[#526174]"
				>
					<CalendarDays size={17} /> Day
				</a>
			</div>

			<div class="space-y-3">
				{#each schedule as item (item.id)}
					{@const Icon = iconFor(item.kind)}
					<div
						class="grid grid-cols-[3.5rem_minmax(0,1fr)] gap-3 sm:grid-cols-[4.5rem_minmax(0,1fr)]"
					>
						<div class="pt-5 text-right">
							<p class="text-sm font-black text-[#243b5c]">{item.time}</p>
							<span
								class={[
									'mt-2 ml-auto block h-2.5 w-2.5 rounded-full',
									item.tone === 'attention'
										? 'bg-[#ee5d4b]'
										: item.tone === 'complete'
											? 'bg-[#4f8a59]'
											: 'bg-[#d3a13c]'
								]}
							></span>
						</div>
						<a
							href={resolve('/households/[slug]', { slug: item.householdId })}
							class="group flex items-center gap-4 rounded-2xl border border-[#e7e1d6] bg-white p-4 shadow-[0_2px_10px_rgba(28,44,67,0.035)] hover:border-[#d8ccb6] hover:shadow-[0_8px_24px_rgba(28,44,67,0.07)] sm:p-5"
						>
							<span
								class={[
									'flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl',
									iconClass(item.tone)
								]}
							>
								<Icon size={22} strokeWidth={2} />
							</span>
							<div class="min-w-0 flex-1">
								<h3 class="truncate text-base font-black sm:text-lg">{item.title}</h3>
								<p class="mt-0.5 truncate text-sm font-semibold text-[#758294]">
									{item.detail} · {item.household}
								</p>
							</div>
							<span
								class={[
									'hidden rounded-full px-3 py-1.5 text-xs font-black sm:block',
									badgeClass(item.tone)
								]}>{item.status}</span
							>
							<ChevronRight
								size={19}
								class="shrink-0 text-[#8793a1] transition-transform group-hover:translate-x-0.5"
							/>
						</a>
					</div>
				{:else}
					<div
						class="rounded-2xl border border-dashed border-[#dcd5c8] bg-[#fffdf8] p-10 text-center"
					>
						<CircleCheck size={32} class="mx-auto text-[#5f9565]" />
						<h3 class="mt-3 font-black">No work scheduled</h3>
						<p class="mt-1 text-sm font-semibold text-[#788596]">
							This day is clear. Add a booking when new work arrives.
						</p>
					</div>
				{/each}
			</div>

			<a
				href={resolve('/calendar')}
				class="mt-5 flex items-center justify-between rounded-2xl border border-[#e7e1d6] bg-[#fffdf8] px-5 py-4 text-sm font-black text-[#243b5c]"
			>
				<span class="flex items-center gap-2"><CalendarDays size={18} /> View full schedule</span>
				<ChevronRight size={18} />
			</a>
		</section>

		<aside class="space-y-5">
			<section
				class="overflow-hidden rounded-2xl border border-[#e7e1d6] bg-white shadow-[0_2px_10px_rgba(28,44,67,0.035)]"
			>
				<div class="flex items-center justify-between border-b border-[#eee8dd] px-5 py-4">
					<h2 class="font-black">Needs attention</h2>
					<span
						class="flex h-7 min-w-7 items-center justify-center rounded-full bg-[#fff0ed] px-2 text-xs font-black text-[#d84832]"
						>{attention.length}</span
					>
				</div>
				<div class="divide-y divide-[#eee8dd]">
					{#each attention.slice(0, 4) as item (item.id)}
						{@const Icon = iconFor(item.kind)}
						<a
							href={resolve('/households/[slug]', { slug: item.householdId })}
							class="flex items-center gap-3 px-5 py-4 hover:bg-[#fffaf3]"
						>
							<span
								class="flex h-10 w-10 items-center justify-center rounded-xl bg-[#fff0ed] text-[#d84832]"
								><Icon size={19} /></span
							>
							<span class="min-w-0 flex-1"
								><span class="block truncate text-sm font-black">{item.title}</span><span
									class="block text-xs font-semibold text-[#7b8797]"
									>{item.time} · {item.household}</span
								></span
							>
							<ChevronRight size={17} class="text-[#8a95a3]" />
						</a>
					{:else}
						<p class="px-5 py-7 text-center text-sm font-semibold text-[#7b8797]">
							Nothing urgent right now.
						</p>
					{/each}
				</div>
			</section>

			<section class="overflow-hidden rounded-2xl border border-[#e7e1d6] bg-white">
				{#each [{ label: 'Today’s work', value: schedule.length, icon: CalendarDays }, { label: 'Completed', value: completedCount, icon: CircleCheck }, { label: 'Boarding', value: boardingCount, icon: Home }] as stat (stat.label)}
					{@const Icon = stat.icon}
					<div class="flex items-center gap-3 border-b border-[#eee8dd] px-5 py-4 last:border-b-0">
						<span
							class="flex h-9 w-9 items-center justify-center rounded-xl bg-[#f8f1e2] text-[#b37b1c]"
							><Icon size={18} /></span
						>
						<span class="flex-1 text-sm font-bold text-[#526174]">{stat.label}</span>
						<span class="text-xl font-black">{stat.value}</span>
					</div>
				{/each}
			</section>

			{#if data.dashboard.openCharges.length > 0}
				<a
					href={resolve('/payments')}
					class="flex items-center gap-3 rounded-2xl border border-[#eadfca] bg-[#fffaf0] p-4"
				>
					<AlertCircle size={20} class="text-[#b37b1c]" />
					<span class="flex-1"
						><span class="block text-sm font-black">Open charges</span><span
							class="block text-xs font-semibold text-[#7b8797]"
							>{data.dashboard.openCharges.length} need review</span
						></span
					>
					<ChevronRight size={18} />
				</a>
			{/if}
		</aside>
	</div>
</main>
