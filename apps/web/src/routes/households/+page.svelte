<script lang="ts">
	import { resolve } from '$app/paths';
	import type { PageData } from './$types';
	import { AlertTriangle, ChevronRight, Dog, Home, Plus, Search } from 'lucide-svelte';

	let { data }: { data: PageData } = $props();

	const households = $derived(
		data.households.map((household) => ({
			name: household.displayName,
			slug: household.id,
			area:
				[household.neighborhood, household.city].filter(Boolean).join(' · ') ||
				'Add an area to simplify routing',
			pets: 'Contacts and pets',
			status: household.active ? 'Active' : 'Inactive',
			next: 'Open profile to add details',
			alert: 'Ready for setup'
		}))
	);
</script>

<svelte:head>
	<title>Households · Pawsear</title>
</svelte:head>

<main class="min-h-screen bg-[#faf9f1] text-prussian-blue-950">
	<div class="mx-auto w-full max-w-4xl px-4 py-5 lg:px-8 lg:py-8">
		<header class="mb-5 flex items-center justify-between gap-3">
			<div>
				<p class="text-xs font-bold tracking-wide text-light-gold-800 uppercase">Households</p>
				<h1 class="text-2xl leading-tight font-black">Casas</h1>
			</div>
			<a
				href={resolve('/households/new')}
				class="flex h-10 w-10 items-center justify-center rounded-lg bg-prussian-blue-700 text-white shadow-sm shadow-black/10"
				aria-label="Add household"
			>
				<Plus size={20} />
			</a>
		</header>

		{#if !data.apiConnected}
			<section
				class="mb-5 rounded-lg border border-sandy-brown-200 bg-sandy-brown-50 p-3 text-sm font-bold text-sandy-brown-800"
			>
				API offline · household data is unavailable
			</section>
		{/if}

		<section
			class="mb-5 rounded-lg border border-light-gold-200 bg-white p-3 shadow-sm shadow-black/5"
		>
			<div class="flex items-center gap-2 rounded-lg bg-[#f6f3e8] px-3 py-2">
				<Search size={18} class="text-light-gold-800" />
				<span class="text-sm font-semibold text-light-gold-800"
					>Search household, pet, or contact</span
				>
			</div>
		</section>

		{#if !data.apiConnected}
			<section
				class="rounded-lg border border-light-gold-200 bg-white p-5 text-center shadow-sm shadow-black/5"
			>
				<h2 class="text-lg font-black">Could not load households</h2>
				<p class="mt-1 text-sm font-semibold text-light-gold-800">
					Start the local API and reload this page.
				</p>
			</section>
		{:else if data.households.length === 0}
			<section
				class="rounded-lg border border-light-gold-200 bg-white p-5 text-center shadow-sm shadow-black/5"
			>
				<div
					class="mx-auto mb-3 flex h-12 w-12 items-center justify-center rounded-lg bg-vanilla-custard-100 text-prussian-blue-800"
				>
					<Home size={23} />
				</div>
				<h2 class="text-lg font-black">No households yet</h2>
				<p class="mx-auto mt-1 max-w-sm text-sm leading-6 font-semibold text-light-gold-800">
					Start with the household, then add its contacts and pets from the profile.
				</p>
				<a
					href={resolve('/households/new')}
					class="mt-4 inline-flex items-center justify-center gap-2 rounded-lg bg-prussian-blue-700 px-4 py-3 text-sm font-black text-white"
				>
					<Plus size={18} />
					Create household
				</a>
			</section>
		{:else}
			<section class="space-y-3">
				{#each households as household (household.slug)}
					<a
						href={resolve('/households/[slug]', { slug: household.slug })}
						class="block rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
					>
						<div class="flex items-start gap-3">
							<div
								class="flex h-11 w-11 shrink-0 items-center justify-center rounded-lg bg-vanilla-custard-100 text-prussian-blue-800"
							>
								<Home size={21} />
							</div>
							<div class="min-w-0 flex-1">
								<div class="flex items-start justify-between gap-3">
									<div class="min-w-0">
										<h2 class="truncate text-lg font-black">{household.name}</h2>
										<p class="text-sm font-semibold text-light-gold-800">{household.area}</p>
									</div>
									<span
										class="rounded-full bg-prussian-blue-100 px-2 py-1 text-xs font-bold text-prussian-blue-800 ring-1 ring-prussian-blue-200"
									>
										{household.status}
									</span>
								</div>

								<div class="mt-3 grid gap-2 text-sm sm:grid-cols-3">
									<p class="flex items-center gap-2 font-bold">
										<Dog size={17} class="text-prussian-blue-700" />
										{household.pets}
									</p>
									<p class="font-semibold text-light-gold-800">{household.next}</p>
									<p class="flex items-center gap-2 font-semibold text-tiger-flame-700">
										<AlertTriangle size={16} />
										{household.alert}
									</p>
								</div>
							</div>
							<ChevronRight size={20} class="mt-2 shrink-0 text-light-gold-800" />
						</div>
					</a>
				{/each}
			</section>
		{/if}
	</div>
</main>
