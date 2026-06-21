<script lang="ts">
	import Navbar from '$lib/components/Navbar.svelte';
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<div class="min-h-screen bg-[#f7f6f2] lg:pl-64">
	<Navbar />
	<div class="pb-20 lg:pb-0">{@render children()}</div>
</div>
