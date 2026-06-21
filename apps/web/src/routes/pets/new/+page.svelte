<script lang="ts">
	import { resolve } from '$app/paths';
	import type { ActionData, PageData } from './$types';
	import { ArrowLeft, Dog, Home, NotebookText, Save } from 'lucide-svelte';

	let { data, form }: { data: PageData; form: ActionData } = $props();
	const fieldValue = (name: keyof NonNullable<ActionData>['values']) => form?.values?.[name] ?? '';
	const selectedHousehold = () => fieldValue('householdId') || data.selectedHouseholdId;
	const selectedHouseholdRecord = $derived(
		data.households.find((household) => household.id === selectedHousehold())
	);
	const canSubmit = $derived(data.households.length > 0);
</script>

<svelte:head>
	<title>New Pet · Pawsear</title>
</svelte:head>

<main class="min-h-screen bg-[#faf9f1] text-prussian-blue-950">
	<div
		class="mx-auto grid w-full max-w-6xl grid-cols-1 lg:grid-cols-[minmax(0,42rem)_20rem] lg:gap-8 lg:px-8"
	>
		<section class="min-w-0">
			<header class="border-b border-light-gold-200 px-4 py-4 lg:border-b-0 lg:px-0 lg:pt-7">
				<a
					href={selectedHousehold()
						? resolve('/households/[slug]', { slug: selectedHousehold() })
						: resolve('/households')}
					class="mb-4 flex items-center gap-2 text-sm font-bold text-light-gold-800"
				>
					<ArrowLeft size={18} />
					Back
				</a>
				<p class="text-xs font-bold tracking-wide text-light-gold-800 uppercase">Pet</p>
				<h1 class="text-2xl leading-tight font-black">Create pet</h1>
			</header>

			<form method="POST" class="space-y-4 px-4 py-4 lg:px-0 lg:pb-10">
				{#if form?.error}
					<section
						class="rounded-lg border border-tiger-flame-200 bg-tiger-flame-50 p-3 text-sm font-bold text-tiger-flame-800"
					>
						{form.error}
					</section>
				{/if}

				{#if data.apiConnected && data.households.length === 0}
					<section
						class="rounded-lg border border-sandy-brown-200 bg-sandy-brown-50 p-3 text-sm font-bold text-sandy-brown-800"
					>
						Create a household before adding pets.
						<a href={resolve('/households/new')} class="ml-1 underline">Create household</a>
					</section>
				{:else if selectedHouseholdRecord}
					<section
						class="rounded-lg border border-prussian-blue-100 bg-prussian-blue-50 p-3 text-sm font-semibold text-prussian-blue-900"
					>
						Adding this pet to <span class="font-black">{selectedHouseholdRecord.displayName}</span
						>.
					</section>
				{/if}

				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
				>
					<div class="mb-4 flex items-center gap-2">
						<Home size={18} class="text-prussian-blue-700" />
						<h2 class="text-base font-black">Household</h2>
					</div>
					<select
						name="householdId"
						required
						class="w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
					>
						<option value="">Choose household</option>
						{#each data.households as household (household.id)}
							<option value={household.id} selected={household.id === selectedHousehold()}
								>{household.displayName}</option
							>
						{/each}
					</select>
				</section>

				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
				>
					<div class="mb-4 flex items-center gap-2">
						<Dog size={18} class="text-prussian-blue-700" />
						<h2 class="text-base font-black">Identity</h2>
					</div>
					<label class="block">
						<span class="text-sm font-black">Name</span>
						<input
							name="name"
							value={fieldValue('name')}
							required
							placeholder="Luna"
							class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
						/>
					</label>
					<div class="mt-3 grid gap-3 sm:grid-cols-2">
						<label class="block">
							<span class="text-sm font-black">Species</span>
							<select
								name="species"
								class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							>
								{#each ['dog', 'cat', 'other'] as species (species)}
									<option value={species} selected={species === (fieldValue('species') || 'dog')}
										>{species}</option
									>
								{/each}
							</select>
						</label>
						<label class="block">
							<span class="text-sm font-black">Breed</span>
							<input
								name="breed"
								value={fieldValue('breed')}
								class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							/>
						</label>
					</div>
					<div class="mt-3 grid gap-3 sm:grid-cols-3">
						<label class="block">
							<span class="text-sm font-black">Size</span>
							<select
								name="size"
								class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							>
								<option value="">Unknown</option>
								{#each ['small', 'medium', 'large', 'giant'] as size (size)}
									<option value={size} selected={size === fieldValue('size')}>{size}</option>
								{/each}
							</select>
						</label>
						<label class="block">
							<span class="text-sm font-black">Sex</span>
							<select
								name="sex"
								class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							>
								<option value="">Unknown</option>
								{#each ['female', 'male'] as sex (sex)}
									<option value={sex} selected={sex === fieldValue('sex')}>{sex}</option>
								{/each}
							</select>
						</label>
						<label class="block">
							<span class="text-sm font-black">Birthdate</span>
							<input
								name="birthdate"
								type="date"
								value={fieldValue('birthdate')}
								class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							/>
						</label>
					</div>
				</section>

				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
				>
					<div class="mb-4 flex items-center gap-2">
						<NotebookText size={18} class="text-prussian-blue-700" />
						<h2 class="text-base font-black">Care notes</h2>
					</div>
					<div class="grid gap-3">
						<textarea
							name="feedingNotes"
							rows="3"
							placeholder="Food schedule and amount"
							class="w-full resize-y rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							>{fieldValue('feedingNotes')}</textarea
						>
						<textarea
							name="medicalNotes"
							rows="3"
							placeholder="Medicine, allergies, chronic conditions"
							class="w-full resize-y rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							>{fieldValue('medicalNotes')}</textarea
						>
						<textarea
							name="behaviorNotes"
							rows="3"
							placeholder="Behavior notes, triggers, leash details"
							class="w-full resize-y rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							>{fieldValue('behaviorNotes')}</textarea
						>
						<textarea
							name="vetNotes"
							rows="3"
							placeholder="Vet or emergency care notes"
							class="w-full resize-y rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							>{fieldValue('vetNotes')}</textarea
						>
						<input
							name="colorMarkings"
							value={fieldValue('colorMarkings')}
							placeholder="Color and markings"
							class="w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
						/>
					</div>
				</section>

				<button
					disabled={!canSubmit}
					class={[
						'flex w-full items-center justify-center gap-2 rounded-lg px-4 py-3 text-sm font-black text-white',
						canSubmit ? 'bg-prussian-blue-700' : 'cursor-not-allowed bg-light-gold-400'
					]}
				>
					<Save size={18} />
					Save pet
				</button>
			</form>
		</section>

		<aside class="hidden border-l border-light-gold-200 bg-[#f6f3e8]/80 px-5 py-7 lg:block">
			<section
				class="sticky top-24 rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
			>
				<h2 class="mb-2 text-base font-black">Care first</h2>
				<p class="text-sm leading-6 font-semibold text-light-gold-800">
					Pet records keep operational notes close to bookings and daily tasks, without assuming a
					single owner or payer.
				</p>
			</section>
		</aside>
	</div>
</main>
