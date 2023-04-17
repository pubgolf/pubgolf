<script>
	import { drawerStore } from '@skeletonlabs/skeleton';
	import SideBarNav from '../nav/SideBarNav.svelte';
	import AboutNav from '../nav/AboutNav.svelte';
</script>

{#if $drawerStore.id === 'admin-nav'}
	<SideBarNav
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
