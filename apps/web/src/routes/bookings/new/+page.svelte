<script lang="ts">
	import { resolve } from '$app/paths';
	import type { ActionData, PageData } from './$types';
	import { untrack } from 'svelte';
	import {
		ArrowLeft,
		CalendarDays,
		Check,
		Clock3,
		Dog,
		Home,
		MessageCircle,
		Route,
		Save
	} from 'lucide-svelte';

	const services = [
		{ value: 'walk', label: 'Walk', description: 'Scheduled walk' },
		{ value: 'visit', label: 'Visit', description: 'Food, water, quick care' },
		{ value: 'client_home_sitting', label: 'Sitting', description: 'Care at client home' },
		{ value: 'boarding', label: 'Boarding', description: 'At caregiver home' },
		{ value: 'transport', label: 'Transport', description: 'Pickup, dropoff, vet run' },
		{ value: 'other', label: 'Other', description: 'Custom service' }
	];

	const statuses = [
		{ value: 'requested', label: 'Requested' },
		{ value: 'confirmed', label: 'Confirmed' }
	];

	const locations = [
		{ value: 'household_home', label: 'Household home' },
		{ value: 'caregiver_home', label: 'Caregiver home' },
		{ value: 'other', label: 'Other' }
	];

	let { data, form }: { data: PageData; form: ActionData } = $props();
	const fieldValue = (name: keyof NonNullable<ActionData>['values']) => form?.values?.[name] ?? '';
	const selectedHouseholdId = $derived(form?.values?.householdId || data.selectedHouseholdId);
	const selectedHousehold = $derived(
		data.households.find((household) => household.id === selectedHouseholdId)
	);
	let selectedService = $state(fieldValue('serviceType') || 'walk');
	let status = $state(fieldValue('status') || 'requested');
	let locationType = $state(fieldValue('locationType') || 'household_home');
	let selectedPetIds = $state<string[]>(untrack(() => form?.values?.petIds ?? data.selectedPetIds));
	let selectedDate = $state(untrack(() => fieldValue('date') || data.defaultDate));
	let selectedStartTime = $state(untrack(() => fieldValue('startTime') || data.defaultStartTime));
	let durationMinutes = $state(fieldValue('durationMinutes') || '45');
	const selectedPets = $derived(data.pets.filter((pet) => selectedPetIds.includes(pet.id)));
	const selectedServiceLabel = $derived(
		services.find((service) => service.value === selectedService)?.label ?? 'Walk'
	);
	const selectedLocationLabel = $derived(
		locations.find((location) => location.value === locationType)?.label ?? 'Household home'
	);
	const canSubmit = $derived(data.apiConnected && data.households.length > 0);

	const changeHousehold = (event: Event) => {
		const select = event.currentTarget as HTMLSelectElement;
		window.location.href = select.value
			? `/bookings/new?householdId=${select.value}`
			: '/bookings/new';
	};
</script>

<svelte:head>
	<title>Create Booking · Pawsear</title>
</svelte:head>

<main class="min-h-screen bg-[#faf9f1] text-prussian-blue-950">
	<div
		class="mx-auto grid w-full max-w-6xl grid-cols-1 lg:grid-cols-[minmax(0,42rem)_20rem] lg:gap-8 lg:px-8"
	>
		<section class="min-w-0">
			<header class="border-b border-light-gold-200 px-4 py-4 lg:border-b-0 lg:px-0 lg:pt-7">
				<a
					href={selectedHouseholdId
						? resolve('/households/[slug]', { slug: selectedHouseholdId })
						: resolve('/')}
					class="mb-4 flex items-center gap-2 text-sm font-bold text-light-gold-800"
				>
					<ArrowLeft size={18} />
					Back
				</a>
				<p class="text-xs font-bold tracking-wide text-light-gold-800 uppercase">Booking</p>
				<h1 class="text-2xl leading-tight font-black">Create booking</h1>
			</header>

			<form id="booking-form" method="POST" class="space-y-4 px-4 py-4 lg:px-0 lg:pb-10">
				<input type="hidden" name="serviceType" value={selectedService} />
				<input type="hidden" name="status" value={status} />
				<input type="hidden" name="locationType" value={locationType} />
				<input type="hidden" name="source" value="manual" />

				{#if form?.error}
					<section
						class="rounded-lg border border-tiger-flame-200 bg-tiger-flame-50 p-3 text-sm font-bold text-tiger-flame-800"
					>
						{form.error}
					</section>
				{/if}

				{#if !data.apiConnected}
					<section
						class="rounded-lg border border-sandy-brown-200 bg-sandy-brown-50 p-3 text-sm font-bold text-sandy-brown-800"
					>
						API offline. Booking creation needs the backend running.
					</section>
				{:else if data.households.length === 0}
					<section
						class="rounded-lg border border-sandy-brown-200 bg-sandy-brown-50 p-3 text-sm font-bold text-sandy-brown-800"
					>
						Create a household before adding bookings.
						<a href={resolve('/households/new')} class="ml-1 underline">Create household</a>
					</section>
				{/if}

				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
				>
					<div class="mb-3 flex items-center gap-2">
						<Route size={18} class="text-prussian-blue-700" />
						<h2 class="text-base font-black">Service</h2>
					</div>
					<div class="grid gap-2 sm:grid-cols-2">
						{#each services as service (service.value)}
							<button
								type="button"
								class={[
									'rounded-lg border p-3 text-left',
									selectedService === service.value
										? 'border-prussian-blue-700 bg-prussian-blue-50 text-prussian-blue-950'
										: 'border-light-gold-200 bg-[#fffdf6]'
								]}
								onclick={() => (selectedService = service.value)}
							>
								<span class="block font-black">{service.label}</span>
								<span class="block text-sm font-semibold text-light-gold-800"
									>{service.description}</span
								>
							</button>
						{/each}
					</div>
				</section>

				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
				>
					<div class="mb-3 flex items-center gap-2">
						<Home size={18} class="text-prussian-blue-700" />
						<h2 class="text-base font-black">Household & pets</h2>
					</div>

					<label class="block">
						<span class="text-sm font-black">Household</span>
						<select
							name="householdId"
							required
							onchange={changeHousehold}
							class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
						>
							<option value="">Choose household</option>
							{#each data.households as household (household.id)}
								<option value={household.id} selected={household.id === selectedHouseholdId}
									>{household.displayName}</option
								>
							{/each}
						</select>
					</label>

					{#if selectedHousehold}
						<div class="mt-3 rounded-lg bg-[#f6f3e8] p-3">
							<p class="text-xs font-bold text-light-gold-800 uppercase">Selected household</p>
							<p class="font-black">{selectedHousehold.displayName}</p>
							<p class="text-sm font-semibold text-light-gold-800">
								{[selectedHousehold.neighborhood, selectedHousehold.city]
									.filter(Boolean)
									.join(' · ') || 'No area yet'}
							</p>
						</div>
					{/if}

					<div class="mt-4">
						<div class="mb-2 flex items-center justify-between gap-3">
							<p class="text-sm font-black">Pets</p>
							{#if selectedHouseholdId}
								<a
									href={resolve(`/pets/new?householdId=${selectedHouseholdId}`)}
									class="text-sm font-bold text-prussian-blue-700">Add pet</a
								>
							{/if}
						</div>
						{#if !selectedHouseholdId}
							<p
								class="rounded-lg border border-light-gold-200 bg-[#fffdf6] p-3 text-sm font-semibold text-light-gold-800"
							>
								Choose a household to load pets.
							</p>
						{:else if data.pets.length === 0}
							<p
								class="rounded-lg border border-light-gold-200 bg-[#fffdf6] p-3 text-sm font-semibold text-light-gold-800"
							>
								No pets yet. This booking will stay at household level unless you add a pet.
							</p>
						{:else}
							<div class="grid grid-cols-2 gap-2">
								{#each data.pets as pet (pet.id)}
									<label
										class={[
											'flex items-center gap-2 rounded-lg border px-3 py-3 text-left font-bold',
											selectedPetIds.includes(pet.id)
												? 'border-prussian-blue-700 bg-prussian-blue-50'
												: 'border-light-gold-200 bg-[#fffdf6]'
										]}
									>
										<input
											class="sr-only"
											type="checkbox"
											name="petIds"
											value={pet.id}
											bind:group={selectedPetIds}
										/>
										<span
											class="flex h-8 w-8 items-center justify-center rounded-lg bg-vanilla-custard-100 text-prussian-blue-800"
										>
											<Dog size={17} />
										</span>
										<span class="min-w-0">
											<span class="block truncate">{pet.name}</span>
											<span class="block truncate text-xs font-semibold text-light-gold-800"
												>{pet.species}</span
											>
										</span>
									</label>
								{/each}
							</div>
						{/if}
					</div>
				</section>

				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
				>
					<div class="mb-3 flex items-center gap-2">
						<CalendarDays size={18} class="text-prussian-blue-700" />
						<h2 class="text-base font-black">Date & time</h2>
					</div>
					<div class="grid gap-3 sm:grid-cols-3">
						<label class="block">
							<span class="text-sm font-black">Date</span>
							<input
								name="date"
								type="date"
								bind:value={selectedDate}
								required
								class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							/>
						</label>
						<label class="block">
							<span class="text-sm font-black">Start</span>
							<input
								name="startTime"
								type="time"
								bind:value={selectedStartTime}
								required
								class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							/>
						</label>
						<label class="block">
							<span class="text-sm font-black">Duration</span>
							<select
								name="durationMinutes"
								bind:value={durationMinutes}
								class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							>
								<option value="30">30 min</option>
								<option value="45">45 min</option>
								<option value="60">60 min</option>
								<option value="90">90 min</option>
								<option value="120">2 hours</option>
								<option value="1440">Full day</option>
							</select>
						</label>
					</div>

					<div class="mt-4 grid grid-cols-2 rounded-lg bg-[#f1ead7] p-1 text-sm font-bold">
						{#each statuses as option (option.value)}
							<button
								type="button"
								class={[
									'rounded-md px-3 py-2',
									status === option.value
										? 'bg-white text-prussian-blue-950 shadow-sm'
										: 'text-light-gold-800'
								]}
								onclick={() => (status = option.value)}
							>
								{option.label}
							</button>
						{/each}
					</div>
				</section>

				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
				>
					<div class="mb-3 flex items-center gap-2">
						<Clock3 size={18} class="text-prussian-blue-700" />
						<h2 class="text-base font-black">Location & notes</h2>
					</div>

					<div class="grid gap-2">
						{#each locations as option (option.value)}
							<button
								type="button"
								class={[
									'flex items-center justify-between rounded-lg border px-3 py-3 text-left font-bold',
									locationType === option.value
										? 'border-prussian-blue-700 bg-prussian-blue-50'
										: 'border-light-gold-200 bg-[#fffdf6]'
								]}
								onclick={() => (locationType = option.value)}
							>
								{option.label}
								{#if locationType === option.value}
									<Check size={18} />
								{/if}
							</button>
						{/each}
					</div>

					<div class="mt-3 grid gap-3 sm:grid-cols-2">
						<label class="block">
							<span class="text-sm font-black">Requested by</span>
							<select
								name="requestedByContactId"
								class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							>
								<option value="">Not set</option>
								{#each data.contacts as contact (`${contact.displayName}-${contact.role}`)}
									<option
										value={contact.contactId}
										selected={contact.contactId === fieldValue('requestedByContactId')}
									>
										{contact.displayName} · {contact.role}
									</option>
								{/each}
							</select>
						</label>
						<label class="block">
							<span class="text-sm font-black">Assigned staff</span>
							<select
								name="assignedStaffId"
								class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							>
								<option value="">Unassigned</option>
								{#each data.staff.filter((member) => member.active) as member (member.id)}
									<option value={member.id} selected={member.id === fieldValue('assignedStaffId')}>
										{member.displayName} · {member.role.replace('_', ' ')}
									</option>
								{/each}
							</select>
						</label>
					</div>

					<label class="mt-3 block">
						<span class="text-sm font-black">Instructions</span>
						<textarea
							name="notes"
							rows="4"
							placeholder="Booking-specific instructions"
							class="mt-2 w-full resize-y rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							>{fieldValue('notes')}</textarea
						>
					</label>
				</section>

				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5 lg:hidden"
				>
					<h2 class="mb-3 text-base font-black">Summary</h2>
					<div class="space-y-2 text-sm font-semibold">
						<p>
							<span class="font-black">{selectedServiceLabel}</span> · {selectedPets.length > 0
								? selectedPets.map((pet) => pet.name).join(' + ')
								: 'Household level'}
						</p>
						<p>
							{selectedHousehold?.displayName ?? 'No household'} · {selectedDate} · {selectedStartTime}
						</p>
						<p>{selectedLocationLabel} · {status}</p>
					</div>
					<button
						disabled={!canSubmit}
						class={[
							'mt-4 flex w-full items-center justify-center gap-2 rounded-lg px-4 py-3 text-sm font-black text-white',
							canSubmit ? 'bg-prussian-blue-700' : 'cursor-not-allowed bg-light-gold-400'
						]}
					>
						<Save size={18} />
						Save booking
					</button>
				</section>
			</form>
		</section>

		<aside class="hidden border-l border-light-gold-200 bg-[#f6f3e8]/80 px-5 py-7 lg:block">
			<div class="sticky top-24 space-y-5">
				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
				>
					<h2 class="mb-3 text-base font-black">Summary</h2>
					<div class="space-y-3 text-sm">
						<div>
							<p class="text-xs font-bold text-light-gold-800 uppercase">Service</p>
							<p class="font-black">{selectedServiceLabel}</p>
						</div>
						<div>
							<p class="text-xs font-bold text-light-gold-800 uppercase">Household</p>
							<p class="font-black">{selectedHousehold?.displayName ?? 'Not selected'}</p>
						</div>
						<div>
							<p class="text-xs font-bold text-light-gold-800 uppercase">Pets</p>
							<p class="font-black">
								{selectedPets.length > 0
									? selectedPets.map((pet) => pet.name).join(' + ')
									: 'Household level'}
							</p>
						</div>
						<div>
							<p class="text-xs font-bold text-light-gold-800 uppercase">Schedule</p>
							<p class="font-black">{selectedDate} · {selectedStartTime} · {durationMinutes} min</p>
						</div>
						<div>
							<p class="text-xs font-bold text-light-gold-800 uppercase">Location</p>
							<p class="font-black">{selectedLocationLabel}</p>
						</div>
						<div>
							<p class="text-xs font-bold text-light-gold-800 uppercase">Status</p>
							<p class="font-black">{status}</p>
						</div>
					</div>
					<button
						form="booking-form"
						disabled={!canSubmit}
						class={[
							'mt-5 flex w-full items-center justify-center gap-2 rounded-lg px-4 py-3 text-sm font-black text-white',
							canSubmit ? 'bg-prussian-blue-700' : 'cursor-not-allowed bg-light-gold-400'
						]}
					>
						<Save size={18} />
						Save booking
					</button>
				</section>

				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 text-sm shadow-sm shadow-black/5"
				>
					<div class="mb-2 flex items-center gap-2">
						<MessageCircle size={18} class="text-prussian-blue-700" />
						<h2 class="text-base font-black">Manual for now</h2>
					</div>
					<p class="leading-6 font-semibold text-light-gold-800">
						Message detection is not connected yet, so this saves as a manual booking with optional
						requested-by context.
					</p>
				</section>
			</div>
		</aside>
	</div>
</main>
