<script>
	import { drawerStore } from '@skeletonlabs/skeleton';
	import SidebarNav from '../nav/SidebarNav.svelte';
	import AboutNav from '../nav/AboutNav.svelte';
</script>

<!-- Touch file -->

{#if $drawerStore.id === 'admin-nav'}
	<SidebarNav
		title="Dashboards"
		items={[
			{
				slug: 'event-info',
				icon: '📝'
			},
			// {
			// 	slug: 'schedule',
			// 	icon: '⏱️'
			// },
			{
				slug: 'players',
				icon: '🏌️'
			},
			{
				slug: 'scores',
				icon: '🏆'
			},
			{
				slug: 'alerts',
				icon: '🚨'
			}
		]}
	/>
{/if}

{#if $drawerStore.id === 'about-nav'}
	<AboutNav />
{/if}
