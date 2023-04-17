<script lang="ts">
	import { Menu } from 'lucide-svelte';
	import { AppShell, AppBar, drawerStore } from '@skeletonlabs/skeleton';
</script>

<AppShell>
	<svelte:fragment slot="pageHeader">
		<AppBar gridColumns="grid-cols-3" slotDefault="place-self-center" shadow="shadow-xl">
			<svelte:fragment slot="lead">
				<button
					class="relative inline-block"
					on:click={() => drawerStore.open({ id: 'admin-nav' })}
				>
					<Menu />
				</button>
			</svelte:fragment>
			<span class="text-l">Admin Panel</span>
		</AppBar>
	</svelte:fragment>
	<div class="container mx-auto p-4">
		<slot />
	</div>
</AppShell>
