<script>
	import '@skeletonlabs/skeleton/themes/theme-skeleton.css';
	import '@skeletonlabs/skeleton/styles/all.css';
	import '../app.postcss';

	import { Drawer, Modal } from '@skeletonlabs/skeleton';
	import SidebarContent from '$lib/components/sidebars/SidebarContent.svelte';
</script>

<Drawer width="w-11/12 max-w-sm">
	<SidebarContent />
</Drawer>

<Modal position="items-start pt-24" />

<div class="h-screen overflow-y-auto">
	<slot />
</div>
