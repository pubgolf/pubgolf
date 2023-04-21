<script lang="ts">
	import { drawerStore } from '@skeletonlabs/skeleton';
	import AboutNav from './AboutNav.svelte';
	import AdminNav from './AdminNav.svelte';
	import AuthFooter from './AuthFooter.svelte';
</script>

{#if $drawerStore.id === 'admin-nav'}
	<AdminNav />

	<AuthFooter />
{/if}

{#if $drawerStore.id === 'about-nav'}
	<AboutNav />
{/if}
