<script>
	import { drawerStore } from '@skeletonlabs/skeleton';

	import AdminNav from './AdminNav.svelte';
</script>

{#if $drawerStore.id === 'admin-nav'}
	<AdminNav />
{/if}
