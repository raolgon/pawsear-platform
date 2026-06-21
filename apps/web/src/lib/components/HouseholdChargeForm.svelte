<script lang="ts">
	type Booking = { id: string; serviceType: string; startAt: string; status: string };
	let { bookings }: { bookings: Booking[] } = $props();
	const completedBookings = $derived(bookings.filter((booking) => booking.status === 'completed'));
</script>

<section class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5">
	<h2 class="text-base font-black">Create charge</h2>
	<p class="mt-1 text-sm font-semibold text-light-gold-800">Record billable work in MXN.</p>
	<form method="POST" action="?/createCharge" class="mt-4 space-y-3">
		<label class="block">
			<span class="text-sm font-black">Booking</span>
			<select
				name="bookingId"
				class="mt-1 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-2"
			>
				<option value="">Household-level charge</option>
				{#each completedBookings as booking (booking.id)}
					<option value={booking.id}
						>{booking.serviceType.replaceAll('_', ' ')} · {booking.startAt.slice(0, 10)} · {booking.status}</option
					>
				{/each}
			</select>
		</label>
		<label class="block">
			<span class="text-sm font-black">Description</span>
			<input
				name="description"
				required
				class="mt-1 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-2"
			/>
		</label>
		<div class="grid gap-3 sm:grid-cols-2">
			<label class="block">
				<span class="text-sm font-black">Amount (MXN)</span>
				<input
					name="amount"
					type="number"
					min="0.01"
					step="0.01"
					required
					class="mt-1 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-2"
				/>
			</label>
			<label class="block">
				<span class="text-sm font-black">Due date</span>
				<input
					name="dueDate"
					type="date"
					class="mt-1 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-2"
				/>
			</label>
		</div>
		<button class="w-full rounded-lg bg-prussian-blue-700 px-4 py-3 text-sm font-black text-white"
			>Create charge</button
		>
	</form>
</section>
