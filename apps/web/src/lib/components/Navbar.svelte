<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import {
		CalendarDays,
		ChevronRight,
		CircleHelp,
		House,
		LayoutDashboard,
		PawPrint,
		ReceiptText
	} from 'lucide-svelte';

	const navItems = [
		{ label: 'Today', href: resolve('/'), icon: LayoutDashboard },
		{ label: 'Calendar', href: resolve('/calendar'), icon: CalendarDays },
		{ label: 'Households', href: resolve('/households'), icon: House },
		{ label: 'Payments', href: resolve('/payments'), icon: ReceiptText }
	];

	const active = (href: string) =>
		href === resolve('/') ? page.url.pathname === href : page.url.pathname.startsWith(href);
</script>

<header class="sticky top-0 z-40 border-b border-[#e9e3d6] bg-[#fffdf8]/95 backdrop-blur lg:hidden">
	<div class="flex h-16 items-center justify-between px-4">
		<a href={resolve('/')} class="flex items-center gap-2 text-[#082652]" aria-label="Pawsear home">
			<span class="flex h-9 w-9 items-center justify-center rounded-xl bg-[#f7ead0] text-[#c58b20]">
				<PawPrint size={21} strokeWidth={2.3} />
			</span>
			<span class="text-xl font-black tracking-[-0.03em]">Pawsear</span>
		</a>
		<span class="rounded-full bg-[#f5f1e8] px-3 py-1.5 text-xs font-extrabold text-[#5f6d80]"
			>Local</span
		>
	</div>
</header>

<aside
	class="fixed inset-y-0 left-0 z-40 hidden w-64 flex-col border-r border-[#e9e3d6] bg-[#fffdf8] lg:flex"
>
	<div class="flex h-24 items-center px-7">
		<a href={resolve('/')} class="flex items-center gap-3 text-[#082652]">
			<span
				class="flex h-11 w-11 items-center justify-center rounded-2xl bg-[#f7ead0] text-[#c58b20]"
			>
				<PawPrint size={25} strokeWidth={2.4} />
			</span>
			<span class="text-2xl font-black tracking-[-0.04em]">Pawsear</span>
		</a>
	</div>

	<nav class="flex-1 space-y-1 px-4 py-5" aria-label="Primary navigation">
		{#each navItems as item (item.href)}
			{@const Icon = item.icon}
			<a
				href={item.href}
				class={[
					'flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-extrabold transition-colors',
					active(item.href)
						? 'bg-[#f5ead4] text-[#082652]'
						: 'text-[#526174] hover:bg-[#f8f4eb] hover:text-[#082652]'
				]}
			>
				<Icon size={20} strokeWidth={2} />
				{item.label}
			</a>
		{/each}
	</nav>

	<div class="space-y-3 p-4">
		<div class="flex items-center gap-3 rounded-2xl border border-[#e9e3d6] bg-white p-3">
			<span
				class="flex h-10 w-10 items-center justify-center rounded-xl bg-[#f7ead0] text-[#c58b20]"
			>
				<PawPrint size={20} />
			</span>
			<div class="min-w-0 flex-1">
				<p class="truncate text-sm font-black text-[#082652]">Pawsear Local</p>
				<p class="text-xs font-semibold text-[#7b8797]">Operations workspace</p>
			</div>
			<ChevronRight size={17} class="text-[#8a95a3]" />
		</div>
		<div class="flex items-center gap-2 px-2 text-xs font-bold text-[#7b8797]">
			<CircleHelp size={16} />
			Local-first MVP
		</div>
	</div>
</aside>

<nav
	class="fixed inset-x-0 bottom-0 z-40 border-t border-[#e5dfd2] bg-[#fffdf8]/97 px-2 pt-2 pb-[max(0.45rem,env(safe-area-inset-bottom))] backdrop-blur lg:hidden"
	aria-label="Mobile navigation"
>
	<div class="grid grid-cols-4">
		{#each navItems as item (item.href)}
			{@const Icon = item.icon}
			<a
				href={item.href}
				class={[
					'flex min-h-14 flex-col items-center justify-center gap-1 rounded-xl text-[0.68rem] font-extrabold',
					active(item.href) ? 'text-[#082652]' : 'text-[#7b8797]'
				]}
			>
				<span
					class={[
						'flex h-8 w-11 items-center justify-center rounded-xl',
						active(item.href) ? 'bg-[#f5ead4]' : ''
					]}
				>
					<Icon size={20} strokeWidth={2.1} />
				</span>
				{item.label}
			</a>
		{/each}
	</div>
</nav>
