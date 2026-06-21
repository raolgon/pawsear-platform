<script lang="ts">
	import { resolve } from '$app/paths';
	import type { ActionData, PageData } from './$types';
	import { ArrowLeft, Home, Phone, Save, UserRound } from 'lucide-svelte';

	let { data, form }: { data: PageData; form: ActionData } = $props();
	const fieldValue = (name: keyof NonNullable<ActionData>['values']) => form?.values?.[name] ?? '';
	const selectedHousehold = () => fieldValue('householdId') || data.selectedHouseholdId;
	const selectedHouseholdRecord = $derived(
		data.households.find((household) => household.id === selectedHousehold())
	);
	const canSubmit = $derived(data.households.length > 0);
</script>

<svelte:head>
	<title>New Contact · Pawsear</title>
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
				<p class="text-xs font-bold tracking-wide text-light-gold-800 uppercase">Contact</p>
				<h1 class="text-2xl leading-tight font-black">Create contact</h1>
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
						Create a household before adding contacts.
						<a href={resolve('/households/new')} class="ml-1 underline">Create household</a>
					</section>
				{:else if selectedHouseholdRecord}
					<section
						class="rounded-lg border border-prussian-blue-100 bg-prussian-blue-50 p-3 text-sm font-semibold text-prussian-blue-900"
					>
						Adding this contact to <span class="font-black"
							>{selectedHouseholdRecord.displayName}</span
						>.
					</section>
				{/if}

				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
				>
					<div class="mb-4 flex items-center gap-2">
						<UserRound size={18} class="text-prussian-blue-700" />
						<h2 class="text-base font-black">Person</h2>
					</div>
					<label class="block">
						<span class="text-sm font-black">Name</span>
						<input
							name="displayName"
							value={fieldValue('displayName')}
							required
							placeholder="Contact name"
							class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
						/>
					</label>
					<div class="mt-3 grid gap-3 sm:grid-cols-2">
						<label class="block">
							<span class="text-sm font-black">Phone</span>
							<input
								name="phone"
								value={fieldValue('phone')}
								placeholder="+52..."
								class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							/>
						</label>
						<label class="block">
							<span class="text-sm font-black">Email</span>
							<input
								name="email"
								value={fieldValue('email')}
								type="email"
								placeholder="name@example.com"
								class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							/>
						</label>
					</div>
					<div class="mt-3 grid gap-3 sm:grid-cols-2">
						<label class="block">
							<span class="text-sm font-black">WhatsApp ID</span>
							<input
								name="whatsappId"
								value={fieldValue('whatsappId')}
								class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							/>
						</label>
						<label class="block">
							<span class="text-sm font-black">Telegram ID</span>
							<input
								name="telegramId"
								value={fieldValue('telegramId')}
								class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							/>
						</label>
					</div>
				</section>

				<section
					class="rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
				>
					<div class="mb-4 flex items-center gap-2">
						<Home size={18} class="text-prussian-blue-700" />
						<h2 class="text-base font-black">Household role</h2>
					</div>
					<label class="block">
						<span class="text-sm font-black">Household</span>
						<select
							name="householdId"
							required
							class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
						>
							<option value="">Choose household</option>
							{#each data.households as household (household.id)}
								<option value={household.id} selected={household.id === selectedHousehold()}
									>{household.displayName}</option
								>
							{/each}
						</select>
					</label>
					<div class="mt-3 grid gap-3 sm:grid-cols-2">
						<label class="block">
							<span class="text-sm font-black">Role</span>
							<select
								name="role"
								class="mt-2 w-full rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
							>
								{#each ['owner', 'partner', 'family', 'domestic_worker', 'payer', 'emergency_contact', 'vet', 'other'] as role (role)}
									<option value={role} selected={role === (fieldValue('role') || 'owner')}
										>{role.replace('_', ' ')}</option
									>
								{/each}
							</select>
						</label>
						<label
							class="flex items-center gap-3 rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-bold"
						>
							<input name="isPrimary" type="checkbox" checked />
							Primary for this role
						</label>
					</div>
					<textarea
						name="notes"
						rows="4"
						placeholder="Contact notes"
						class="mt-3 w-full resize-y rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
						>{fieldValue('notes')}</textarea
					>
					<textarea
						name="relationshipNotes"
						rows="3"
						placeholder="Notes for this household role"
						class="mt-3 w-full resize-y rounded-lg border border-light-gold-200 bg-[#fffdf6] px-3 py-3 font-semibold outline-none focus:border-prussian-blue-700"
						>{fieldValue('relationshipNotes')}</textarea
					>
				</section>

				<button
					disabled={!canSubmit}
					class={[
						'flex w-full items-center justify-center gap-2 rounded-lg px-4 py-3 text-sm font-black text-white',
						canSubmit ? 'bg-prussian-blue-700' : 'cursor-not-allowed bg-light-gold-400'
					]}
				>
					<Save size={18} />
					Save contact
				</button>
			</form>
		</section>

		<aside class="hidden border-l border-light-gold-200 bg-[#f6f3e8]/80 px-5 py-7 lg:block">
			<section
				class="sticky top-24 rounded-lg border border-light-gold-200 bg-white p-4 shadow-sm shadow-black/5"
			>
				<div class="mb-2 flex items-center gap-2">
					<Phone size={18} class="text-prussian-blue-700" />
					<h2 class="text-base font-black">Why contacts are separate</h2>
				</div>
				<p class="text-sm leading-6 font-semibold text-light-gold-800">
					A contact can schedule, pay, hand off a pet, or act as emergency contact. Roles belong to
					the household relationship.
				</p>
			</section>
		</aside>
	</div>
</main>
