<script lang="ts">
	import { resolve } from '$app/paths';
	import type { ActionData } from './$types';
	import { ArrowLeft, Home, MapPin, NotebookText, Save } from 'lucide-svelte';

	let { form }: { form: ActionData } = $props();

	const fieldValue = (name: keyof NonNullable<ActionData>['values']) => form?.values?.[name] ?? '';
</script>

<svelte:head>
	<title>New Household · Pawsear</title>
</svelte:head>

<main class="min-h-screen bg-[#faf9f1] text-prussian-blue-950">
	<div
		class="mx-auto grid w-full max-w-6xl grid-cols-1 lg:grid-cols-[minmax(0,42rem)_20rem] lg:gap-8 lg:px-8"
	>
		<section class="min-w-0">
			<header class="border-b border-light-gold-200 px-4 py-4 lg:border-b-0 lg:px-0 lg:pt-7">
				<div class="mb-4 flex items-center justify-between">
					<a
						href={resolve('/households')}
						class="flex items-center gap-2 text-sm font-bold text-light-gold-800"
					>
						<ArrowLeft size={18} />
						Back
					</a>
				</div>
				<p class="text-xs font-bold tracking-wide text-light-gold-800 uppercase">Household</p>
				<h1 class="text-2xl leading-tight font-black">Create household</h1>
			</header>

			<form method="POST" class="space-y-4 px-4 py-4 lg:px-0 lg:pb-10">
				{#if form?.error}
					<section
						class="rounded-lg border border-tiger-flame-200 bg-tiger-flame-50 p-3 text-sm font-bold text-tiger-flame-800"
					>
						{form.error}
					</section>
				{/if}

				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
				>
					<div class="mb-4 flex items-center gap-2">
						<Home size={18} class="text-prussian-blue-700" />
						<h2 class="text-base font-black">Identity</h2>
					</div>
					<label class="block">
						<span class="text-sm font-black">Household name</span>
						<input
							name="displayName"
							value={fieldValue('displayName')}
							required
							placeholder="Casa de Luna y Max"
							class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 text-base font-semibold outline-none focus:border-prussian-blue-700"
						/>
					</label>
				</section>

				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
				>
					<div class="mb-4 flex items-center gap-2">
						<MapPin size={18} class="text-prussian-blue-700" />
						<h2 class="text-base font-black">Address</h2>
					</div>
					<div class="grid gap-3">
						<label class="block">
							<span class="text-sm font-black">Address line 1</span>
							<input
								name="addressLine1"
								value={fieldValue('addressLine1')}
								placeholder="Street and number"
								class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							/>
						</label>
						<label class="block">
							<span class="text-sm font-black">Address line 2</span>
							<input
								name="addressLine2"
								value={fieldValue('addressLine2')}
								placeholder="Apartment, gate code, references"
								class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							/>
						</label>
						<div class="grid gap-3 sm:grid-cols-2">
							<label class="block">
								<span class="text-sm font-black">Neighborhood</span>
								<input
									name="neighborhood"
									value={fieldValue('neighborhood')}
									placeholder="Roma Norte"
									class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
								/>
							</label>
							<label class="block">
								<span class="text-sm font-black">City</span>
								<input
									name="city"
									value={fieldValue('city') || 'CDMX'}
									class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
								/>
							</label>
						</div>
						<label class="block">
							<span class="text-sm font-black">Timezone</span>
							<input
								name="timezone"
								value={fieldValue('timezone') || 'America/Mexico_City'}
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
						<h2 class="text-base font-black">Operational notes</h2>
					</div>
					<textarea
						name="notes"
						rows="5"
						placeholder="Access notes, preferred handoff, recurring reminders..."
						class="w-full resize-y rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
						>{fieldValue('notes')}</textarea
					>
				</section>

				<div
					class="sticky bottom-0 -mx-4 border-t border-light-gold-200 bg-[#faf9f1]/95 px-4 py-3 backdrop-blur lg:static lg:mx-0 lg:border-0 lg:bg-transparent lg:px-0"
				>
					<button
						class="flex w-full items-center justify-center gap-2 rounded-lg bg-prussian-blue-700 px-4 py-3 text-sm font-black text-white shadow-sm shadow-black/10"
					>
						<Save size={18} />
						Save household
					</button>
				</div>
			</form>
		</section>

		<aside class="hidden border-l border-light-gold-200 bg-[#f6f3e8]/80 px-5 py-7 lg:block">
			<div class="sticky top-24 space-y-4">
				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
				>
					<h2 class="mb-2 text-base font-black">What happens next?</h2>
					<p class="text-sm leading-6 font-semibold text-light-gold-800">
						After creating the household, add contacts and pets from its profile. Keeping household
						first preserves the real-world relationship between owners, payers, pets, and staff.
					</p>
				</section>
				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
				>
					<h2 class="mb-2 text-base font-black">MVP rule</h2>
					<p class="text-sm leading-6 font-semibold text-light-gold-800">
						Clients do not need an account. This record is for internal operations, scheduling, care
						notes, and payment tracking.
					</p>
				</section>
			</div>
		</aside>
	</div>
</main>
