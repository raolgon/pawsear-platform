<script lang="ts">
	import PaymentForm from '$lib/components/PaymentForm.svelte';
	import type { ActionData, PageData } from './$types';
	import { AlertTriangle, ChevronDown, ChevronRight, ReceiptText, Wallet } from 'lucide-svelte';

	type Tab = 'charges' | 'payments';
	let { data, form }: { data: PageData; form: ActionData } = $props();
	let activeTab = $state<Tab>('charges');

	const money = (amountMinor: number, currency = 'MXN') =>
		new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency,
			maximumFractionDigits: 0
		}).format(amountMinor / 100);

	const openTotal = $derived(
		data.charges.reduce((total, charge) => total + charge.outstandingMinor, 0)
	);
	const receivedToday = $derived(
		data.payments
			.filter(
				(payment) => payment.receivedAt.slice(0, 10) === new Date().toISOString().slice(0, 10)
			)
			.reduce((total, payment) => total + payment.amountMinor, 0)
	);
	const groupedCharges = $derived(
		Object.entries(
			data.charges.reduce<Record<string, typeof data.charges>>((groups, charge) => {
				(groups[charge.householdId] ??= []).push(charge);
				return groups;
			}, {})
		)
	);
</script>

<svelte:head><title>Payments · Pawsear</title></svelte:head>

<main class="min-h-screen bg-[#f7f6f2] text-[#071b3b]">
	<header class="border-b border-[#e9e3d6] bg-[#fffdf8]">
		<div class="mx-auto max-w-[1360px] px-4 py-6 sm:px-7 lg:px-10">
			<p class="text-sm font-bold text-[#7b8797]">Money</p>
			<h1 class="mt-0.5 text-3xl font-black tracking-[-0.04em] lg:text-4xl">Payments</h1>
		</div>
	</header>

	<div
		class="mx-auto grid max-w-[1360px] gap-6 px-4 py-6 sm:px-7 lg:grid-cols-[minmax(0,1fr)_minmax(390px,0.8fr)] lg:px-10 lg:py-8"
	>
		<section class="min-w-0 space-y-5">
			{#if !data.apiConnected}
				<div
					class="rounded-2xl border border-[#f3c7bf] bg-[#fff4f1] p-4 text-sm font-bold text-[#a63e2f]"
				>
					API offline · payment data is unavailable
				</div>
			{/if}
			{#if form?.error}
				<div
					class="rounded-2xl border border-[#f3c7bf] bg-[#fff4f1] p-4 text-sm font-bold text-[#a63e2f]"
				>
					{form.error}
				</div>
			{/if}

			<div class="grid grid-cols-2 gap-3">
				<div class="rounded-2xl border border-[#e7e1d6] bg-white p-4 sm:p-5">
					<p class="text-xs font-black tracking-[0.1em] text-[#8a95a3] uppercase">Outstanding</p>
					<p class="mt-2 text-2xl font-black tracking-[-0.03em]">{money(openTotal)}</p>
					<p class="mt-1 text-xs font-semibold text-[#7b8797]">
						{data.charges.length} open charges
					</p>
				</div>
				<div class="rounded-2xl border border-[#e7e1d6] bg-white p-4 sm:p-5">
					<p class="text-xs font-black tracking-[0.1em] text-[#8a95a3] uppercase">Received today</p>
					<p class="mt-2 text-2xl font-black tracking-[-0.03em] text-[#397546]">
						{money(receivedToday)}
					</p>
					<p class="mt-1 text-xs font-semibold text-[#7b8797]">Manual payments</p>
				</div>
			</div>

			<div class="flex rounded-xl bg-[#ece8df] p-1 lg:hidden">
				<button
					class={[
						'flex-1 rounded-lg px-3 py-2 text-sm font-black',
						activeTab === 'charges' ? 'bg-white shadow-sm' : 'text-[#687587]'
					]}
					onclick={() => (activeTab = 'charges')}>Open charges</button
				>
				<button
					class={[
						'flex-1 rounded-lg px-3 py-2 text-sm font-black',
						activeTab === 'payments' ? 'bg-white shadow-sm' : 'text-[#687587]'
					]}
					onclick={() => (activeTab = 'payments')}>History</button
				>
			</div>

			<section class={activeTab === 'payments' ? 'hidden lg:block' : 'block'}>
				<div
					class="overflow-hidden rounded-2xl border border-[#e7e1d6] bg-white shadow-[0_2px_10px_rgba(28,44,67,0.035)]"
				>
					<div class="flex items-center justify-between border-b border-[#ece6db] px-5 py-4">
						<div>
							<h2 class="text-lg font-black">Open charges</h2>
							<p class="mt-0.5 text-xs font-semibold text-[#7b8797]">
								Select allocations in the payment form.
							</p>
						</div>
						<span
							class="flex h-8 min-w-8 items-center justify-center rounded-full bg-[#f5ead4] px-2 text-xs font-black text-[#9a6714]"
							>{data.charges.length}</span
						>
					</div>

					{#each groupedCharges as [householdId, householdCharges] (householdId)}
						<div class="border-b border-[#ece6db] last:border-b-0">
							<div class="flex items-center gap-3 bg-[#fcfaf5] px-5 py-3">
								<span
									class="flex h-9 w-9 items-center justify-center rounded-xl bg-[#f7ead0] text-[#b57d1c]"
									><ReceiptText size={18} /></span
								>
								<div class="min-w-0 flex-1">
									<p class="truncate text-sm font-black">
										{data.householdNames[householdId] ?? 'Add household details'}
									</p>
									<p class="text-xs font-semibold text-[#7b8797]">
										{householdCharges.length}
										{householdCharges.length === 1 ? 'charge' : 'charges'}
									</p>
								</div>
								<p class="font-black">
									{money(
										householdCharges.reduce((total, charge) => total + charge.outstandingMinor, 0),
										householdCharges[0]?.currency
									)}
								</p>
								<ChevronDown size={17} class="text-[#85909d]" />
							</div>
							<div class="divide-y divide-[#f0ebe2]">
								{#each householdCharges as charge (charge.id)}
									<div class="grid grid-cols-[minmax(0,1fr)_auto] gap-3 px-5 py-4 pl-10 sm:pl-16">
										<div class="min-w-0">
											<p class="truncate text-sm font-black">{charge.description}</p>
											<p class="mt-1 text-xs font-semibold text-[#7b8797]">
												{charge.dueDate ? `Due ${charge.dueDate.slice(0, 10)}` : 'Add a due date'} · {charge.status.replaceAll(
													'_',
													' '
												)}
											</p>
										</div>
										<p class="font-black tabular-nums">
											{money(charge.outstandingMinor, charge.currency)}
										</p>
									</div>
								{/each}
							</div>
						</div>
					{:else}
						<div class="p-10 text-center">
							<ReceiptText size={30} class="mx-auto text-[#6b9870]" />
							<h3 class="mt-3 font-black">No open charges</h3>
							<p class="mt-1 text-sm font-semibold text-[#7b8797]">Everything recorded is paid.</p>
						</div>
					{/each}
				</div>
			</section>

			<section class={activeTab === 'charges' ? 'hidden lg:block' : 'block'}>
				<div class="overflow-hidden rounded-2xl border border-[#e7e1d6] bg-white">
					<div class="border-b border-[#ece6db] px-5 py-4">
						<h2 class="text-lg font-black">Recent payments</h2>
					</div>
					<div class="divide-y divide-[#ece6db]">
						{#each data.payments.slice(0, 8) as payment (payment.id)}
							<div class="flex items-center gap-3 px-5 py-4">
								<span
									class="flex h-10 w-10 items-center justify-center rounded-xl bg-[#eaf5e6] text-[#397546]"
									><Wallet size={19} /></span
								>
								<div class="min-w-0 flex-1">
									<p class="truncate text-sm font-black">
										{payment.payerContactId
											? (data.contactNames[payment.payerContactId] ?? 'Add payer details')
											: 'Add who made this payment'}
									</p>
									<p class="text-xs font-semibold text-[#7b8797]">
										{payment.method.replaceAll('_', ' ')} · {payment.receivedAt.slice(0, 10)}
									</p>
								</div>
								<p class="font-black text-[#397546]">
									{money(payment.amountMinor, payment.currency)}
								</p>
								<ChevronRight size={17} class="text-[#85909d]" />
							</div>
						{:else}
							<p class="p-8 text-center text-sm font-semibold text-[#7b8797]">
								Payments will appear here after you record the first one.
							</p>
						{/each}
					</div>
				</div>
			</section>
		</section>

		<aside class="order-first lg:order-last">
			<div class="lg:sticky lg:top-6">
				<PaymentForm
					contacts={data.contacts}
					charges={data.charges}
					householdNames={data.householdNames}
				/>
				<div
					class="mt-3 flex items-start gap-2 rounded-xl px-3 py-2 text-xs font-semibold text-[#7b8797]"
				>
					<AlertTriangle size={16} class="mt-0.5 shrink-0" />
					This records a manual payment only. No card processing.
				</div>
			</div>
		</aside>
	</div>
</main>
