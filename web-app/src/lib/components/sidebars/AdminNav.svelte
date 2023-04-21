<script lang="ts">
	import SidebarNav from './SidebarNav.svelte';
</script>

<SidebarNav
	title="Dashboards"
	items={[
		// {
		// 	slug: 'event-info',
		// 	icon: '📝'
		// },
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
		}
		// {
		// 	slug: 'alerts',
		// 	icon: '🚨'
		// }
	]}
/>
