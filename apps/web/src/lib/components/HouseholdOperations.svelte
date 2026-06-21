<script lang="ts">
	type Booking = {
		id: string;
		serviceType: string;
		status: string;
		startAt: string;
	};

	type CareTask = {
		id: string;
		title: string;
		taskType: string;
		status: string;
		dueAt?: string | null;
	};

	let { bookings, careTasks }: { bookings: Booking[]; careTasks: CareTask[] } = $props();

	const nextBookingStatus = (status: string) => {
		if (status === 'requested') return { value: 'confirmed', label: 'Confirm' };
		if (status === 'confirmed') return { value: 'in_progress', label: 'Start' };
		if (status === 'in_progress') return { value: 'completed', label: 'Complete' };
		return null;
	};

	const timeLabel = (value?: string | null) =>
		value?.slice(0, 16).replace('T', ' ') ?? 'Set a time for this task';
</script>

<section class="space-y-3">
	<div class="flex items-center justify-between">
		<h2 class="text-base font-black">Bookings</h2>
		<span class="text-xs font-bold text-light-gold-800">{bookings.length}</span>
	</div>
	{#each bookings as booking (booking.id)}
		{@const next = nextBookingStatus(booking.status)}
		<article class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5">
			<div class="flex items-start justify-between gap-3">
				<div>
					<p class="font-black capitalize">{booking.serviceType.replaceAll('_', ' ')}</p>
					<p class="text-sm font-semibold text-light-gold-800">{timeLabel(booking.startAt)}</p>
				</div>
				<span
					class="rounded-full bg-prussian-blue-100 px-2 py-1 text-xs font-bold text-prussian-blue-800"
				>
					{booking.status.replaceAll('_', ' ')}
				</span>
			</div>
			{#if next || !['completed', 'cancelled'].includes(booking.status)}
				<div class="mt-3 flex justify-end gap-2">
					{#if !['completed', 'cancelled'].includes(booking.status)}
						<form method="POST" action="?/bookingStatus">
							<input type="hidden" name="bookingId" value={booking.id} />
							<input type="hidden" name="status" value="cancelled" />
							<button class="rounded-lg border border-light-gold-200 px-3 py-2 text-sm font-bold"
								>Cancel</button
							>
						</form>
					{/if}
					{#if next}
						<form method="POST" action="?/bookingStatus">
							<input type="hidden" name="bookingId" value={booking.id} />
							<input type="hidden" name="status" value={next.value} />
							<button class="rounded-lg bg-prussian-blue-700 px-3 py-2 text-sm font-bold text-white"
								>{next.label}</button
							>
						</form>
					{/if}
				</div>
			{/if}
		</article>
	{:else}
		<p
			class="rounded-lg border border-light-gold-200 bg-white p-4 text-sm font-semibold text-light-gold-800"
		>
			Create the first booking to start this household's schedule.
		</p>
	{/each}
</section>

<section class="space-y-3">
	<div class="flex items-center justify-between">
		<h2 class="text-base font-black">Care tasks</h2>
		<span class="text-xs font-bold text-light-gold-800">{careTasks.length}</span>
	</div>
	{#each careTasks as task (task.id)}
		<article class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5">
			<div class="flex items-start justify-between gap-3">
				<div>
					<p class="font-black">{task.title}</p>
					<p class="text-sm font-semibold text-light-gold-800">
						{task.taskType.replaceAll('_', ' ')} · {timeLabel(task.dueAt)}
					</p>
				</div>
				<span class="text-xs font-bold capitalize">{task.status}</span>
			</div>
			{#if task.status === 'pending'}
				<div class="mt-3 grid gap-2 sm:grid-cols-[auto_minmax(0,1fr)]">
					<form method="POST" action="?/careTaskStatus">
						<input type="hidden" name="taskId" value={task.id} />
						<input type="hidden" name="status" value="completed" />
						<button
							class="w-full rounded-lg bg-prussian-blue-700 px-3 py-2 text-sm font-bold text-white"
							>Complete</button
						>
					</form>
					<form method="POST" action="?/careTaskStatus" class="flex min-w-0 gap-2">
						<input type="hidden" name="taskId" value={task.id} />
						<input type="hidden" name="status" value="skipped" />
						<input
							name="skippedReason"
							required
							placeholder="Reason to skip"
							class="min-w-0 flex-1 rounded-lg border border-light-gold-200 px-3 py-2 text-sm"
						/>
						<button class="rounded-lg border border-light-gold-200 px-3 py-2 text-sm font-bold"
							>Skip</button
						>
					</form>
				</div>
			{/if}
		</article>
	{:else}
		<p
			class="rounded-lg border border-light-gold-200 bg-white p-4 text-sm font-semibold text-light-gold-800"
		>
			Care tasks will appear here when this household needs timed care.
		</p>
	{/each}
</section>
