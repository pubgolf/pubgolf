<script>
	import '@skeletonlabs/skeleton/themes/theme-skeleton.css';
	import '@skeletonlabs/skeleton/styles/all.css';
	import '../app.postcss';

	import { Drawer, Modal, Toast } from '@skeletonlabs/skeleton';
	import SidebarContent from '$lib/components/sidebars/Sidebar.svelte';
</script>

<Drawer width="w-11/12 max-w-sm">
	<SidebarContent />
</Drawer>

<Modal position="items-start pt-24" />

<Toast />

<div class="h-screen overflow-y-auto">
	<slot />
</div>
