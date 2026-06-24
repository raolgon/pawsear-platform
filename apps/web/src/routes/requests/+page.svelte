<script lang="ts">
	import { resolve } from '$app/paths';
	import type { ActionData, PageData } from './$types';
	import { AlertTriangle, ArrowRight, Inbox, MessageSquareText, Plus, Send } from 'lucide-svelte';
	import { SvelteURLSearchParams } from 'svelte/reactivity';

	let { data, form }: { data: PageData; form: ActionData } = $props();

	const titleCase = (value: string) =>
		value.replaceAll('_', ' ').replace(/\b\w/g, (letter) => letter.toUpperCase());
	const pendingStatuses = ['needs_review', 'needs_more_info', 'confirmed'];
	const pendingRequests = $derived(
		data.detectedRequests.filter((request) => pendingStatuses.includes(request.status))
	);
	const historyRequests = $derived(
		data.detectedRequests.filter((request) =>
			['converted_to_booking', 'ignored'].includes(request.status)
		)
	);
	const visibleRequests = $derived(
		data.view === 'history'
			? historyRequests
			: data.view === 'all'
				? data.detectedRequests
				: pendingRequests
	);
	const bookingReviewHref = (request: PageData['detectedRequests'][number]) => {
		const params = new SvelteURLSearchParams({ requestId: request.id });
		if (request.householdId) params.set('householdId', request.householdId);
		if (request.startAt) params.set('date', request.startAt.slice(0, 10));
		return `/bookings/new?${params.toString()}`;
	};
	const latestReplyFor = (requestId: string, templateKey?: string) =>
		data.outboundMessages.find(
			(message) =>
				message.detectedRequestId === requestId &&
				(!templateKey || message.templateKey === templateKey)
		);
	const replyLabel = (templateKey: string) =>
		({
			request_details: 'Details requested',
			booking_confirmed: 'Confirmation',
			request_declined: 'Decline notice'
		})[templateKey] ?? 'Telegram reply';
	const householdOptionsFor = (contactId: string) =>
		data.contactHouseholds.filter((link) => link.contactId === contactId);
</script>

<svelte:head><title>Request review · Pawsear</title></svelte:head>

<main class="mx-auto w-full max-w-6xl px-4 py-6 sm:px-6 lg:px-8 lg:py-8">
	<header class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
		<div>
			<p class="text-xs font-black tracking-[0.18em] text-[#a8741b] uppercase">Message capture</p>
			<h1 class="mt-1 text-3xl font-black tracking-[-0.04em] text-[#082652]">Request review</h1>
			<p class="mt-2 max-w-2xl text-sm font-semibold text-[#687789]">
				Messages stay here until a person confirms what should become scheduled work.
			</p>
		</div>
		<span class="rounded-full bg-[#f5ead4] px-3 py-1.5 text-xs font-black text-[#805b19]">
			{pendingRequests.length} pending
		</span>
	</header>
	<nav class="mb-5 flex gap-2" aria-label="Request views">
		<a
			href={resolve('/requests?view=pending')}
			class={[
				'rounded-xl px-4 py-2 text-sm font-black',
				data.view === 'pending' ? 'bg-[#082d60] text-white' : 'border border-[#ddd7cc] bg-white'
			]}>Pending ({pendingRequests.length})</a
		>
		<a
			href={resolve('/requests?view=history')}
			class={[
				'rounded-xl px-4 py-2 text-sm font-black',
				data.view === 'history' ? 'bg-[#082d60] text-white' : 'border border-[#ddd7cc] bg-white'
			]}>History ({historyRequests.length})</a
		>
		<a
			href={resolve('/requests?view=all')}
			class={[
				'rounded-xl px-4 py-2 text-sm font-black',
				data.view === 'all' ? 'bg-[#082d60] text-white' : 'border border-[#ddd7cc] bg-white'
			]}>All ({data.detectedRequests.length})</a
		>
	</nav>

	{#if !data.apiConnected}
		<div
			class="mb-5 flex gap-3 rounded-2xl border border-[#ecc7bd] bg-[#fff3ef] p-4 text-[#8b3f31]"
		>
			<AlertTriangle size={20} class="mt-0.5 shrink-0" />
			<div>
				<p class="font-black">The local API is offline</p>
				<p class="text-sm font-semibold">{data.apiError}</p>
			</div>
		</div>
	{/if}
	{#if form?.error}
		<p
			class="mb-5 rounded-2xl border border-[#ecc7bd] bg-[#fff3ef] p-4 text-sm font-bold text-[#8b3f31]"
		>
			{form.error}
		</p>
	{:else if form?.success}
		<p
			class="mb-5 rounded-2xl border border-[#c9dfc7] bg-[#f0f8ed] p-4 text-sm font-bold text-[#397546]"
		>
			{form.success}
		</p>
	{/if}

	<div class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_22rem]">
		<section class="space-y-4">
			{#each visibleRequests as request (request.id)}
				{@const latestReply = latestReplyFor(request.id)}
				{@const confirmationReply = latestReplyFor(request.id, 'booking_confirmed')}
				<article class="overflow-hidden rounded-2xl border border-[#e7e1d6] bg-white">
					<div class="flex flex-wrap items-center gap-2 border-b border-[#eee8dd] px-5 py-3">
						<span class="rounded-full bg-[#edf4ff] px-2.5 py-1 text-xs font-black text-[#1c5c98]"
							>{titleCase(request.channel)}</span
						>
						<span class="rounded-full bg-[#fff4df] px-2.5 py-1 text-xs font-black text-[#8a6118]"
							>{titleCase(request.status)}</span
						>
						<span class="ml-auto text-xs font-bold text-[#7b8797]"
							>Confidence: {request.confidence}</span
						>
					</div>
					<div class="p-5">
						<p class="text-[0.95rem] leading-6 font-semibold whitespace-pre-wrap text-[#263d5c]">
							{request.body}
						</p>
						<div class="mt-4 grid gap-2 text-sm sm:grid-cols-2">
							<p class="rounded-xl bg-[#f8f5ee] px-3 py-2">
								<span class="font-black">Household:</span>
								{request.householdName || 'Match a household during review'}
							</p>
							<p class="rounded-xl bg-[#f8f5ee] px-3 py-2">
								<span class="font-black">Contact:</span>
								{request.contactName || 'Match the sender during review'}
							</p>
							<p class="rounded-xl bg-[#f8f5ee] px-3 py-2">
								<span class="font-black">Service:</span>
								{request.serviceType
									? titleCase(request.serviceType)
									: 'Choose the requested service'}
							</p>
							<p class="rounded-xl bg-[#f8f5ee] px-3 py-2">
								<span class="font-black">When:</span>
								{request.startAt ?? 'Confirm the requested date and time'}
							</p>
						</div>
						{#if latestReply}
							<div
								class={[
									'mt-4 rounded-xl border px-3 py-2 text-sm font-bold',
									latestReply.status === 'sent'
										? 'border-[#c9dfc7] bg-[#f0f8ed] text-[#397546]'
										: latestReply.status === 'failed'
											? 'border-[#ecc7bd] bg-[#fff3ef] text-[#8b3f31]'
											: 'border-[#bfd5ef] bg-[#edf4ff] text-[#1c5c98]'
								]}
							>
								{replyLabel(latestReply.templateKey)} · {titleCase(latestReply.status)}
								{#if latestReply.lastError}
									<span class="block text-xs">{latestReply.lastError}</span>
								{/if}
							</div>
						{/if}
						<div class="mt-4 flex flex-wrap gap-2">
							{#if request.senderExternalId && !request.contactId}
								<form
									method="POST"
									action="?/linkContact"
									class="flex w-full flex-col gap-2 rounded-xl border border-[#bfd5ef] bg-[#edf4ff] p-3 sm:flex-row"
								>
									<input type="hidden" name="id" value={request.id} />
									{#if data.contacts.length > 0}
										<label class="min-w-0 flex-1">
											<span class="sr-only">Message sender</span>
											<select
												name="contactId"
												required
												class="w-full rounded-lg border border-[#aac5e3] bg-white px-3 py-2.5 text-sm font-bold text-[#263d5c]"
											>
												<option value="">Who sent this message?</option>
												{#each data.contacts as contact (contact.id)}
													<option value={contact.id}>{contact.displayName}</option>
												{/each}
											</select>
										</label>
									{:else}
										<label class="min-w-0 flex-1">
											<span class="sr-only">Sender name</span>
											<input
												name="displayName"
												required
												placeholder="Sender name"
												class="w-full rounded-lg border border-[#aac5e3] bg-white px-3 py-2.5 text-sm font-bold text-[#263d5c]"
											/>
										</label>
										<label class="min-w-0 flex-1">
											<span class="sr-only">Sender household</span>
											<select
												name="householdId"
												required
												class="w-full rounded-lg border border-[#aac5e3] bg-white px-3 py-2.5 text-sm font-bold text-[#263d5c]"
											>
												<option value="">Choose sender household</option>
												{#each data.households as household (household.id)}
													<option value={household.id}>{household.displayName}</option>
												{/each}
											</select>
										</label>
									{/if}
									<button class="rounded-lg bg-[#1c5c98] px-4 py-2.5 text-sm font-black text-white">
										Remember sender
									</button>
								</form>
							{:else if request.contactId && !request.householdId}
								{#if householdOptionsFor(request.contactId).length === 0}
									<p
										class="w-full rounded-xl border border-[#e6cf9e] bg-[#fff8e8] px-3 py-2 text-sm font-bold text-[#604819]"
									>
										{request.contactName} is recognized, but must be linked to a household before it can
										be selected.
									</p>
								{:else}
									<form
										method="POST"
										action="?/linkHousehold"
										class="flex w-full flex-col gap-2 rounded-xl border border-[#e6cf9e] bg-[#fff8e8] p-3 sm:flex-row"
									>
										<input type="hidden" name="id" value={request.id} />
										<label class="min-w-0 flex-1">
											<span class="sr-only">Household for this conversation</span>
											<select
												name="householdId"
												required
												class="w-full rounded-lg border border-[#dfc27e] bg-white px-3 py-2.5 text-sm font-bold text-[#604819]"
											>
												<option value="">Which household is this about?</option>
												{#each householdOptionsFor(request.contactId) as household (household.householdId)}
													<option value={household.householdId}>{household.householdName}</option>
												{/each}
											</select>
										</label>
										<button
											class="rounded-lg bg-[#8a6118] px-4 py-2.5 text-sm font-black text-white"
										>
											Remember household
										</button>
									</form>
								{/if}
							{:else if request.contactId && request.householdId}
								<p
									class="w-full rounded-xl bg-[#f0f8ed] px-3 py-2 text-sm font-bold text-[#397546]"
								>
									Future {titleCase(request.channel)} messages in this chat will recognize {request.contactName}
									and {request.householdName}.
								</p>
							{/if}
							{#if request.status === 'converted_to_booking'}
								<a
									href={resolve(
										`/calendar?date=${(request.convertedBookingStartAt ?? request.startAt ?? request.createdAt).slice(0, 10)}`
									)}
									class="inline-flex items-center gap-2 rounded-xl bg-[#397546] px-4 py-2.5 text-sm font-black text-white"
									>View in agenda <ArrowRight size={16} /></a
								>
								{#if !confirmationReply || confirmationReply.status === 'failed'}
									<form method="POST" action="?/queueReply">
										<input type="hidden" name="id" value={request.id} />
										<input type="hidden" name="templateKey" value="booking_confirmed" />
										<button
											class="inline-flex items-center gap-2 rounded-xl border border-[#9fc7a5] px-4 py-2.5 text-sm font-black text-[#397546]"
										>
											<Send size={16} />
											{confirmationReply ? 'Retry confirmation' : 'Send confirmation'}
										</button>
									</form>
								{/if}
							{:else if request.status === 'ignored'}
								<form method="POST" action="?/updateStatus">
									<input type="hidden" name="id" value={request.id} />
									<input type="hidden" name="status" value="needs_review" />
									<button
										class="rounded-xl border border-[#d9d3c8] px-4 py-2.5 text-sm font-black text-[#42546b]"
									>
										Reopen request
									</button>
								</form>
							{:else}
								<form method="GET" action={bookingReviewHref(request)}>
									<button
										class="inline-flex items-center gap-2 rounded-xl bg-[#082d60] px-4 py-2.5 text-sm font-black text-white"
										>Review booking <ArrowRight size={16} /></button
									>
								</form>
								{#if latestReply?.status !== 'pending'}
									<form method="POST" action="?/queueReply">
										<input type="hidden" name="id" value={request.id} />
										<input type="hidden" name="templateKey" value="request_details" />
										<button
											class="inline-flex items-center gap-2 rounded-xl border border-[#d9d3c8] px-4 py-2.5 text-sm font-black text-[#42546b]"
											><Send size={16} /> Ask via Telegram</button
										>
									</form>
									<form method="POST" action="?/queueReply">
										<input type="hidden" name="id" value={request.id} />
										<input type="hidden" name="templateKey" value="request_declined" />
										<button
											class="rounded-xl border border-[#ecc7bd] px-4 py-2.5 text-sm font-black text-[#8b3f31]"
										>
											Decline & notify
										</button>
									</form>
								{/if}
								<form method="POST" action="?/updateStatus">
									<input type="hidden" name="id" value={request.id} />
									<input type="hidden" name="status" value="ignored" />
									<button class="rounded-xl px-4 py-2.5 text-sm font-black text-[#7b8797]">
										Ignore
									</button>
								</form>
							{/if}
						</div>
					</div>
				</article>
			{:else}
				<div
					class="rounded-2xl border border-dashed border-[#dcd5c8] bg-[#fffdf8] p-10 text-center"
				>
					<Inbox size={34} class="mx-auto text-[#a8741b]" />
					<h2 class="mt-3 text-lg font-black">
						{data.view === 'history'
							? 'No reviewed requests yet'
							: 'No messages waiting for review'}
					</h2>
					<p class="mt-1 text-sm font-semibold text-[#758294]">
						Paste a client message or connect n8n to start the queue.
					</p>
				</div>
			{/each}
		</section>

		<aside>
			<form
				method="POST"
				action="?/importManual"
				class="rounded-2xl border border-[#e7e1d6] bg-white p-5 lg:sticky lg:top-6"
			>
				<div class="flex items-center gap-3">
					<span
						class="flex h-10 w-10 items-center justify-center rounded-xl bg-[#f5ead4] text-[#a8741b]"
						><MessageSquareText size={20} /></span
					>
					<div>
						<h2 class="font-black">Paste a message</h2>
						<p class="text-xs font-semibold text-[#758294]">WhatsApp, Telegram, or a call note</p>
					</div>
				</div>
				<label class="mt-4 block"
					><span class="text-sm font-black">Original message</span><textarea
						name="body"
						rows="8"
						required
						placeholder="Paste the client's exact message here…"
						class="mt-2 w-full resize-y rounded-xl border border-[#ddd7cc] bg-[#fffdf8] px-3 py-3 text-sm font-semibold outline-none focus:border-[#1c5c98]"
					></textarea></label
				>
				<p class="mt-2 text-xs leading-5 font-semibold text-[#7b8797]">
					Pawsear preserves this text as context. It will not create a booking automatically.
				</p>
				<button
					class="mt-4 flex w-full items-center justify-center gap-2 rounded-xl bg-[#082d60] px-4 py-3 text-sm font-black text-white"
					><Plus size={17} /> Add to review queue</button
				>
				<div
					class="mt-4 flex items-start gap-2 rounded-xl bg-[#f7f4ec] p-3 text-xs font-semibold text-[#687789]"
				>
					<Send size={15} class="mt-0.5 shrink-0" />n8n will use the same import endpoint for
					automated capture.
				</div>
			</form>
		</aside>
	</div>
</main>
