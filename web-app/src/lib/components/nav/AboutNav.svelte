<script lang="ts">
	import SideBarNav from './SideBarNav.svelte';
</script>

<SideBarNav
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
