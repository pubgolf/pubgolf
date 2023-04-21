<script lang="ts">
	import { AppBar, AppShell, drawerStore } from '@skeletonlabs/skeleton';
	import { Menu } from 'lucide-svelte';
	import AboutNav from '$lib/components/sidebars/AboutNav.svelte';
	import SiteFooter from '$lib/components/SiteFooter.svelte';
</script>

<AppShell>
	<svelte:fragment slot="sidebarLeft">
		<div class="hidden md:block">
			<AboutNav />
		</div>
	</svelte:fragment>

	<svelte:fragment slot="pageHeader">
		<div class="md:hidden">
			<AppBar gridColumns="grid-cols-3" slotDefault="place-self-center" shadow="shadow-xl">
				<svelte:fragment slot="lead">
					<button
						class="relative inline-block"
						on:click={() => drawerStore.open({ id: 'about-nav' })}
					>
						<Menu />
					</button>
				</svelte:fragment>
				<span class="text-l">About PubGolf</span>
			</AppBar>
		</div>
	</svelte:fragment>

	<div class="container mx-auto p-4 max-w-prose">
		<slot />
	</div>

	<svelte:fragment slot="pageFooter">
		<SiteFooter />
	</svelte:fragment>
</AppShell>
