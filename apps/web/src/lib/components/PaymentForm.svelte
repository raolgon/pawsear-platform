<script lang="ts">
	import { SvelteDate } from 'svelte/reactivity';
	type Contact = { id: string; displayName: string };
	type Charge = {
		id: string;
		householdId: string;
		description: string;
		currency: string;
		outstandingMinor: number;
	};

	let {
		contacts,
		charges,
		householdNames
	}: { contacts: Contact[]; charges: Charge[]; householdNames: Record<string, string> } = $props();

	const money = (amountMinor: number, currency: string) =>
		new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency,
			maximumFractionDigits: 0
		}).format(amountMinor / 100);
	const localNow = () => {
		const date = new SvelteDate();
		date.setMinutes(date.getMinutes() - date.getTimezoneOffset());
		return date.toISOString().slice(0, 16);
	};
</script>

<section class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5">
	<h2 class="text-base font-black">Record payment</h2>
	<form method="POST" action="?/recordPayment" class="mt-4 space-y-3">
		<label class="block">
			<span class="text-sm font-black">Payer</span>
			<select
				name="payerContactId"
				class="mt-1 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-2"
			>
				<option value="">Unknown payer</option>
				{#each contacts as contact (contact.id)}
					<option value={contact.id}>{contact.displayName}</option>
				{/each}
			</select>
		</label>
		<div class="grid gap-3 sm:grid-cols-2">
			<label class="block">
				<span class="text-sm font-black">Amount</span>
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
				<span class="text-sm font-black">Received</span>
				<input
					name="receivedAt"
					type="datetime-local"
					value={localNow()}
					required
					class="mt-1 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-2"
				/>
			</label>
		</div>
		<label class="block">
			<span class="text-sm font-black">Method</span>
			<select
				name="method"
				class="mt-1 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-2"
			>
				<option value="bank_transfer">Transfer</option>
				<option value="cash">Cash</option>
				<option value="card_external">External card</option>
				<option value="other">Other</option>
			</select>
		</label>
		<label class="block">
			<span class="text-sm font-black">Reference</span>
			<input
				name="reference"
				class="mt-1 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-2"
			/>
		</label>

		<div>
			<p class="text-sm font-black">Allocate to charges</p>
			<p class="text-xs font-semibold text-light-gold-800">
				Select one or more charges and enter the amount for each.
			</p>
			<div class="mt-2 space-y-2">
				{#each charges as charge (charge.id)}
					<div class="rounded-lg border border-light-gold-200 p-3">
						<label class="flex items-start gap-2">
							<input type="checkbox" name="chargeId" value={charge.id} class="mt-1" />
							<span class="min-w-0 flex-1">
								<span class="block font-bold">{charge.description}</span>
								<span class="block text-xs font-semibold text-light-gold-800"
									>{householdNames[charge.householdId] ?? 'Unknown household'} · {money(
										charge.outstandingMinor,
										charge.currency
									)} outstanding</span
								>
							</span>
						</label>
						<input
							name={`allocation:${charge.id}`}
							type="number"
							min="0.01"
							max={charge.outstandingMinor / 100}
							step="0.01"
							placeholder="Allocation"
							class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-2"
						/>
					</div>
				{:else}
					<p class="rounded-lg bg-[#f6f3e8] p-3 text-sm font-semibold text-light-gold-800">
						No open charges to allocate.
					</p>
				{/each}
			</div>
		</div>

		<label class="block">
			<span class="text-sm font-black">Notes</span>
			<textarea
				name="notes"
				rows="2"
				class="mt-1 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-2"
			></textarea>
		</label>
		<button class="w-full rounded-lg bg-prussian-blue-700 px-4 py-3 text-sm font-black text-white"
			>Save payment</button
		>
	</form>
</section>
