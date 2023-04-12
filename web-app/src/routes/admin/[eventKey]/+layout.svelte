<script lang="ts">
	import { Menu } from 'lucide-svelte';
	import { AppShell, AppBar, Drawer, drawerStore } from '@skeletonlabs/skeleton';
	import { page } from '$app/stores';

	const dashboards = ['scores', 'schedule', 'players', 'alerts'] as const;

	function titleCase(s: string) {
		return s.charAt(0).toUpperCase() + s.slice(1);
	}
</script>

<Drawer width="w-11/12 max-w-sm">
	<nav class="list-nav">
		<span>Nav</span>
		<ul>
			{#each dashboards as dashboardName}
				<li>
					<a
						class:bg-primary-500={$page.route.id?.endsWith(dashboardName)}
						on:click={() => drawerStore.close()}
						href="../{dashboardName}/"
					>
						<span class="badge bg-primary-500">💀</span>
						<span class="flex-auto">{titleCase(dashboardName)}</span>
					</a>
				</li>
			{/each}
		</ul>
	</nav>
</Drawer>

<AppShell>
	<svelte:fragment slot="pageHeader">
		<AppBar gridColumns="grid-cols-3" slotDefault="place-self-center">
			<svelte:fragment slot="lead">
				<button
					class="relative inline-block"
					on:click={() => drawerStore.open({ id: 'admin-nav' })}
				>
					<Menu />
				</button>
			</svelte:fragment>
			Admin Panel
		</AppBar>
	</svelte:fragment>
	<slot />
</AppShell>
