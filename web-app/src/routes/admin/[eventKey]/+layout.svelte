<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { USER_NOT_AUTHORIZED_ERROR } from '$lib/auth/client';
	import type { DisplayError } from '$lib/components/ErrorBanner.svelte';
	import ErrorBanner from '$lib/components/ErrorBanner.svelte';
	import SiteFooter from '$lib/components/SiteFooter.svelte';
	import { getAdminServiceClient } from '$lib/rpc/client';
	import { AppBar, AppShell, drawerStore } from '@skeletonlabs/skeleton';
	import { Menu, RefreshCwIcon } from 'lucide-svelte';
	import { onMount } from 'svelte';

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
	<svelte:fragment slot="header">
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
			<div class="card py-12 flex flex-col items-center">
				<p class="mb-4">Checking auth status...</p>
				<RefreshCwIcon class="animate-spin" />
			</div>
		{/if}
	</div>
	<svelte:fragment slot="pageFooter">
		<SiteFooter />
	</svelte:fragment>
</AppShell>
