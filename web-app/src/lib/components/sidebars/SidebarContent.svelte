<script>
	import { drawerStore } from '@skeletonlabs/skeleton';

	import AboutNav from './AboutNav.svelte';
	import AdminNav from './AdminNav.svelte';
</script>

{#if $drawerStore.id === 'admin-nav'}
	<AdminNav />
{/if}

{#if $drawerStore.id === 'about-nav'}
	<AboutNav />
{/if}
