<script lang="ts">
	import { resolve } from '$app/paths';
	import HouseholdChargeForm from '$lib/components/HouseholdChargeForm.svelte';
	import HouseholdOperations from '$lib/components/HouseholdOperations.svelte';
	import type { ActionData, PageData } from './$types';
	import {
		ArrowLeft,
		ChevronRight,
		CircleCheck,
		Dog,
		House,
		MapPin,
		MessageCircle,
		Phone,
		Plus,
		ReceiptText,
		Trash2,
		Utensils
	} from 'lucide-svelte';

	let { data, form }: { data: PageData; form: ActionData } = $props();
</script>

<svelte:head><title>{data.household.name} · Pawsear</title></svelte:head>

<main class="min-h-screen bg-[#f7f6f2] text-[#071b3b]">
	<header class="border-b border-[#e9e3d6] bg-[#fffdf8]">
		<div class="mx-auto max-w-[1180px] px-4 py-5 sm:px-7 lg:px-10 lg:py-7">
			<div class="flex items-center justify-between">
				<a
					href={resolve('/households')}
					class="flex h-10 w-10 items-center justify-center rounded-xl border border-[#e2ddd2] bg-white text-[#24405f]"
					aria-label="Back to households"
				>
					<ArrowLeft size={19} />
				</a>
				<a
					href={resolve(`/bookings/new?householdId=${data.household.id}`)}
					class="flex h-10 w-10 items-center justify-center rounded-xl bg-[#082d60] text-white shadow-sm"
					aria-label="Create booking"
				>
					<Plus size={20} />
				</a>
			</div>

			<div class="mt-6 flex items-center gap-4">
				<span
					class="flex h-16 w-16 shrink-0 items-center justify-center rounded-3xl bg-[#f7ead0] text-[#b77d1b] sm:h-20 sm:w-20"
				>
					<House size={31} strokeWidth={1.8} />
				</span>
				<div class="min-w-0">
					<p class="text-xs font-black tracking-[0.13em] text-[#8a95a3] uppercase">Household</p>
					<h1 class="mt-1 truncate text-3xl font-black tracking-[-0.04em] sm:text-4xl">
						{data.household.name}
					</h1>
					<p class="mt-1 flex items-center gap-1.5 text-sm font-bold text-[#758294]">
						<MapPin size={16} />
						{data.household.area}
					</p>
				</div>
			</div>
		</div>
	</header>

	<div
		class="mx-auto grid max-w-[1180px] gap-6 px-4 py-6 sm:px-7 lg:grid-cols-[minmax(0,1fr)_310px] lg:px-10 lg:py-8"
	>
		<section class="min-w-0 space-y-6">
			{#if form?.error}
				<div
					class="rounded-2xl border border-[#f3c7bf] bg-[#fff4f1] p-4 text-sm font-bold text-[#a63e2f]"
				>
					{form.error}
				</div>
			{/if}

			<section class="flex items-center gap-3 rounded-2xl border border-[#d8e5d5] bg-[#fbfdf9] p-4">
				<span
					class="flex h-11 w-11 shrink-0 items-center justify-center rounded-2xl bg-[#dcefd8] text-[#397546]"
					><CircleCheck size={23} /></span
				>
				<div class="min-w-0 flex-1">
					<p class="font-black">Household up to date</p>
					<p class="text-sm font-semibold text-[#758294]">{data.household.alerts}</p>
				</div>
				<ChevronRight size={19} class="text-[#85909d]" />
			</section>

			<section>
				<div class="mb-3 flex items-center justify-between">
					<div>
						<p class="text-xs font-black tracking-[0.12em] text-[#8a95a3] uppercase">
							Care profiles
						</p>
						<h2 class="mt-1 text-xl font-black">Pets</h2>
					</div>
					<a
						href={resolve(`/pets/new?householdId=${data.household.id}`)}
						class="rounded-xl border border-[#e2ddd2] bg-white px-3 py-2 text-sm font-black text-[#183c68]"
						>Add pet</a
					>
				</div>
				<div class="grid gap-3 sm:grid-cols-2">
					{#each data.household.pets as pet (pet.id ?? pet.name)}
						<article
							class="rounded-2xl border border-[#e7e1d6] bg-white p-4 shadow-[0_2px_10px_rgba(28,44,67,0.035)]"
						>
							<div class="flex items-center gap-3">
								<span
									class="flex h-12 w-12 items-center justify-center rounded-2xl bg-[#f8f1e2] text-[#a8741b]"
									><Dog size={23} /></span
								>
								<div class="min-w-0 flex-1">
									<h3 class="truncate text-xl font-black">{pet.name}</h3>
									<p class="text-sm font-semibold text-[#758294]">{pet.kind}</p>
								</div>
								<ChevronRight size={18} class="text-[#85909d]" />
							</div>
							<div class="mt-4 space-y-2 border-t border-[#eee8dd] pt-3 text-sm">
								<p class="flex gap-2 font-bold">
									<Utensils size={17} class="mt-0.5 shrink-0 text-[#a8741b]" /><span
										>{pet.care}</span
									>
								</p>
								<p class="text-[#657386]">
									<span class="font-black text-[#243b5c]">Note:</span>
									{pet.note}
								</p>
							</div>
							<a
								href={resolve(
									`/bookings/new?householdId=${data.household.id}${pet.id ? `&petId=${pet.id}` : ''}`
								)}
								class="mt-4 flex w-full items-center justify-center rounded-xl bg-[#082d60] px-3 py-2.5 text-sm font-black text-white"
								>Book care</a
							>
						</article>
					{:else}
						<div
							class="sm:col-span-2 rounded-2xl border border-dashed border-[#dcd5c8] bg-[#fffdf8] p-8 text-center"
						>
							<Dog size={30} class="mx-auto text-[#a8741b]" />
							<h3 class="mt-3 font-black">No pets yet</h3>
							<p class="mt-1 text-sm font-semibold text-[#758294]">
								Add a pet to store care instructions.
							</p>
						</div>
					{/each}
				</div>
			</section>

			<section class="rounded-2xl border border-[#e7e1d6] bg-white p-5">
				<div class="flex items-center justify-between">
					<h2 class="text-lg font-black">Next booking</h2>
					<span class="rounded-full bg-[#edf4ff] px-3 py-1 text-xs font-black text-[#1c5c98]"
						>{data.household.nextWork.status}</span
					>
				</div>
				<div class="mt-4 flex items-center gap-4">
					<span
						class="flex h-12 w-12 items-center justify-center rounded-2xl bg-[#f8f1e2] text-[#a8741b]"
						><ReceiptText size={22} /></span
					>
					<div class="min-w-0 flex-1">
						<p class="text-xl font-black">{data.household.nextWork.time}</p>
						<p class="text-sm font-semibold text-[#758294]">{data.household.nextWork.pets}</p>
					</div>
					<a
						href={resolve(`/bookings/new?householdId=${data.household.id}`)}
						class="hidden rounded-xl border border-[#d6deea] px-4 py-2 text-sm font-black text-[#183c68] sm:block"
						>New booking</a
					>
				</div>
			</section>

			<HouseholdOperations bookings={data.bookings} careTasks={data.careTasks} />

			<section class="rounded-2xl border border-[#e7e1d6] bg-white p-5">
				<div class="flex items-center justify-between">
					<h2 class="text-lg font-black">Contacts</h2>
					<a
						href={resolve(`/contacts/new?householdId=${data.household.id}`)}
						class="text-sm font-black text-[#1c5c98]">Add</a
					>
				</div>
				<div class="mt-3 divide-y divide-[#eee8dd]">
					{#each data.household.contacts as contact (`${contact.name}-${contact.role}`)}
						<div class="flex items-center gap-3 py-3">
							<span
								class="flex h-10 w-10 items-center justify-center rounded-xl bg-[#f8f1e2] text-[#526174]"
								><House size={18} /></span
							>
							<div class="min-w-0 flex-1">
								<p class="truncate font-black">{contact.name}</p>
								<p class="truncate text-xs font-semibold text-[#758294]">{contact.role}</p>
							</div>
							<button
								class="flex h-9 w-9 items-center justify-center rounded-full border border-[#e2ddd2]"
								aria-label={`Call ${contact.name}`}><Phone size={16} /></button
							>
							<button
								class="flex h-9 w-9 items-center justify-center rounded-full border border-[#e2ddd2]"
								aria-label={`Message ${contact.name}`}><MessageCircle size={16} /></button
							>
						</div>
					{:else}
						<p class="py-5 text-sm font-semibold text-[#758294]">
							Add the first contact for scheduling, handoffs, or payments.
						</p>
					{/each}
				</div>
			</section>

			<div class="lg:hidden"><HouseholdChargeForm bookings={data.bookings} /></div>
		</section>

		<aside class="space-y-5">
			<a
				href={resolve('/payments')}
				class="flex items-center gap-3 rounded-2xl border border-[#dfd1b4] bg-[#fffaf0] p-4"
			>
				<span
					class="flex h-11 w-11 items-center justify-center rounded-2xl bg-[#f7ead0] text-[#a8741b]"
					><ReceiptText size={21} /></span
				>
				<div class="min-w-0 flex-1">
					<p class="text-xs font-black tracking-[0.1em] text-[#8a95a3] uppercase">Balance</p>
					<p class="mt-1 text-lg font-black">{data.household.balance}</p>
				</div>
				<ChevronRight size={19} />
			</a>
			<div class="hidden lg:block"><HouseholdChargeForm bookings={data.bookings} /></div>
			<section class="rounded-2xl border border-[#e7e1d6] bg-white p-4">
				<h2 class="font-black">Household notes</h2>
				<p class="mt-2 text-sm leading-6 font-semibold text-[#657386]">{data.household.notes}</p>
			</section>
			<section class="rounded-2xl border border-[#edc7c1] bg-[#fff8f6] p-4">
				<div class="flex items-center gap-2 text-[#9f3f32]">
					<Trash2 size={18} />
					<h2 class="font-black">Delete household</h2>
				</div>
				<p class="mt-2 text-sm font-semibold text-[#7d5b56]">
					Permanently deletes its pets, bookings, tasks, charges, and message history. This cannot
					be undone.
				</p>
				<form method="POST" action="?/deleteHousehold" class="mt-3 space-y-2">
					<label class="block">
						<span class="text-xs font-black text-[#7d5b56]">
							Type “{data.household.name}” to confirm
						</span>
						<input
							name="confirmationName"
							required
							autocomplete="off"
							class="mt-1 w-full rounded-xl border border-[#dfb4ad] bg-white px-3 py-2.5 text-sm font-bold"
						/>
					</label>
					<button class="w-full rounded-xl bg-[#a63e2f] px-3 py-2.5 text-sm font-black text-white">
						Delete permanently
					</button>
				</form>
			</section>
		</aside>
	</div>
</main>
