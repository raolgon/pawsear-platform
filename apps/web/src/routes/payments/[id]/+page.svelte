<script lang="ts">
	import { resolve } from '$app/paths';
	import type { ActionData, PageData } from './$types';
	import {
		ArrowLeft,
		CircleCheck,
		Download,
		FileImage,
		FileText,
		ReceiptText
	} from 'lucide-svelte';

	let { data, form }: { data: PageData; form: ActionData } = $props();
	const money = (amountMinor: number, currency: string) =>
		new Intl.NumberFormat('es-MX', {
			style: 'currency',
			currency,
			minimumFractionDigits: 2
		}).format(amountMinor / 100);
	const methodLabel: Record<string, string> = {
		cash: 'Efectivo',
		bank_transfer: 'Transferencia bancaria',
		card_external: 'Tarjeta externa',
		other: 'Otro método'
	};
</script>

<svelte:head><title>Payment receipt · Pawsear</title></svelte:head>

<main class="min-h-screen bg-[#f7f6f2] text-[#071b3b]">
	<header class="border-b border-[#e9e3d6] bg-[#fffdf8]">
		<div class="mx-auto max-w-[1180px] px-4 py-5 sm:px-7 lg:px-10">
			<a
				href={resolve('/payments')}
				class="inline-flex items-center gap-2 text-sm font-black text-[#476079]"
			>
				<ArrowLeft size={18} /> Payments
			</a>
			<div class="mt-5 flex items-center gap-4">
				<span
					class="flex h-14 w-14 items-center justify-center rounded-2xl bg-[#eaf5e6] text-[#397546]"
				>
					<ReceiptText size={27} />
				</span>
				<div>
					<p class="text-xs font-black tracking-[0.12em] text-[#8a95a3] uppercase">
						Payment received
					</p>
					<h1 class="mt-1 text-3xl font-black tracking-[-0.04em]">
						{money(data.payment.amountMinor, data.payment.currency)}
					</h1>
				</div>
			</div>
		</div>
	</header>

	<div
		class="mx-auto grid max-w-[1180px] gap-6 px-4 py-6 sm:px-7 lg:grid-cols-[minmax(0,1fr)_330px] lg:px-10 lg:py-8"
	>
		<section class="min-w-0">
			{#if form?.error}
				<div
					class="mb-5 rounded-2xl border border-[#f3c7bf] bg-[#fff4f1] p-4 text-sm font-bold text-[#a63e2f]"
				>
					{form.error}
				</div>
			{/if}

			{#if data.receipt}
				<div
					class="mb-4 flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-[#c9dfc7] bg-[#f0f8ed] p-4"
				>
					<div class="flex items-center gap-3">
						<CircleCheck size={22} class="text-[#397546]" />
						<div>
							<p class="font-black text-[#285f35]">Receipt ready</p>
							<p class="text-xs font-bold text-[#5e7b64]">{data.receipt.receiptNumber}</p>
						</div>
					</div>
					<div class="flex gap-2">
						<a
							href={resolve(`/payments/${data.payment.id}/receipt.pdf?download=1`)}
							class="inline-flex items-center gap-2 rounded-xl bg-[#082d60] px-3 py-2 text-sm font-black text-white"
						>
							<FileText size={17} /> PDF <Download size={15} />
						</a>
						<a
							href={resolve(`/payments/${data.payment.id}/receipt.png?download=1`)}
							class="inline-flex items-center gap-2 rounded-xl border border-[#c4d3c2] bg-white px-3 py-2 text-sm font-black text-[#285f35]"
						>
							<FileImage size={17} /> PNG <Download size={15} />
						</a>
					</div>
				</div>
				<div
					class="overflow-hidden rounded-2xl border border-[#e3ddd2] bg-[#eeeae1] p-3 shadow-sm sm:p-6"
				>
					<img
						src={resolve('/payments/[id]/receipt.png', { id: data.payment.id })}
						alt={`Receipt ${data.receipt.receiptNumber}`}
						class="mx-auto h-auto w-full max-w-[700px] rounded-xl shadow-[0_12px_45px_rgba(20,35,55,0.16)]"
					/>
				</div>
			{:else}
				<div class="rounded-2xl border border-[#e4ded2] bg-white p-7 text-center shadow-sm">
					<span
						class="mx-auto flex h-16 w-16 items-center justify-center rounded-3xl bg-[#f5ead4] text-[#9a6714]"
					>
						<ReceiptText size={30} />
					</span>
					<h2 class="mt-4 text-xl font-black">Create the payment receipt</h2>
					<p class="mx-auto mt-2 max-w-md text-sm leading-6 font-semibold text-[#758294]">
						This freezes the current payer, amount, method, reference, and applied charges. It will
						not change the payment.
					</p>
					<form method="POST" action="?/issueReceipt" class="mt-5">
						<button class="rounded-xl bg-[#082d60] px-5 py-3 text-sm font-black text-white"
							>Issue receipt</button
						>
					</form>
				</div>
			{/if}
		</section>

		<aside class="space-y-4 lg:sticky lg:top-6 lg:self-start">
			<section class="rounded-2xl border border-[#e4ded2] bg-white p-5 shadow-sm">
				<h2 class="font-black">Payment details</h2>
				<dl class="mt-4 space-y-4 text-sm">
					<div>
						<dt class="text-xs font-black text-[#8a95a3] uppercase">Payer</dt>
						<dd class="mt-1 font-bold">{data.payerName}</dd>
					</div>
					<div>
						<dt class="text-xs font-black text-[#8a95a3] uppercase">Received</dt>
						<dd class="mt-1 font-bold">
							{new Date(data.payment.receivedAt).toLocaleString('es-MX')}
						</dd>
					</div>
					<div>
						<dt class="text-xs font-black text-[#8a95a3] uppercase">Method</dt>
						<dd class="mt-1 font-bold">{methodLabel[data.payment.method] ?? 'Otro método'}</dd>
					</div>
					<div>
						<dt class="text-xs font-black text-[#8a95a3] uppercase">Reference</dt>
						<dd class="mt-1 font-bold">{data.payment.reference ?? 'No reference recorded'}</dd>
					</div>
					<div>
						<dt class="text-xs font-black text-[#8a95a3] uppercase">Applied charges</dt>
						<dd class="mt-1 font-bold">{data.payment.allocations.length}</dd>
					</div>
				</dl>
			</section>
			<p class="rounded-xl bg-[#f1eee6] p-3 text-xs leading-5 font-semibold text-[#758294]">
				Internal payment confirmation only. It is not an invoice or CFDI.
			</p>
		</aside>
	</div>
</main>
