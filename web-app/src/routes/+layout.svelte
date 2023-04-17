<script>
	import '@skeletonlabs/skeleton/themes/theme-skeleton.css';
	import '@skeletonlabs/skeleton/styles/all.css';
	import '../app.postcss';

	import { Drawer, Modal } from '@skeletonlabs/skeleton';
	import SidebarContent from '$lib/components/sidebars/SideBarContent.svelte';
</script>

<Drawer width="w-11/12 max-w-sm">
	<SidebarContent />
</Drawer>

<Modal />

<div class="h-screen overflow-y-auto">
	<slot />
</div>
