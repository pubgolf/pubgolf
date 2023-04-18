<script lang="ts">
	import SidebarNav from './SidebarNav.svelte';
</script>

<SidebarNav
	title="About PubGolf"
	items={[
		{
			title: 'Privacy Policy',
			slug: 'privacy',
			icon: '🕵️'
		},
		{
			title: 'Contact Us',
			slug: 'contact',
			icon: '📞'
		}
	]}
/>
