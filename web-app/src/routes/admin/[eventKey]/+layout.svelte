<script lang="ts">
	import { Menu } from 'lucide-svelte';
	import { AppShell, AppBar, drawerStore } from '@skeletonlabs/skeleton';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { USER_NOT_AUTHORIZED_ERROR } from '$lib/auth/client';
	import { page } from '$app/stores';
	import type { DisplayError } from '$lib/components/ErrorBanner.svelte';
	import ErrorBanner from '$lib/components/ErrorBanner.svelte';
	import { getAdminServiceClient } from '$lib/rpc/client';

	let authInit = false;
	let authError: DisplayError = null;
	function retryAuth() {
		authError = null;
		setupAuth();
	}

	async function setupAuth() {
		try {
			await getAdminServiceClient();
			authInit = true;
		} catch (error) {
			if (error === USER_NOT_AUTHORIZED_ERROR) {
				goto(`/auth/login?redirect=${$page.url.pathname}`);
				return;
			}
			authError = error as DisplayError;
		}
	}

	onMount(setupAuth);
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
		<ErrorBanner error={authError} dismissLabel="Retry" on:dismiss={retryAuth} />
		{#if authInit}
			<slot />
		{:else}
			Authenticating...
		{/if}
	</div>
	<svelte:fragment slot="pageFooter">
		<SiteFooter />
	</svelte:fragment>
</AppShell>
