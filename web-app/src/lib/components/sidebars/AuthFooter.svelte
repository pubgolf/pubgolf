<script lang="ts">
	import { drawerStore, toastStore } from '@skeletonlabs/skeleton';
	import { logOut } from '$lib/auth/client';
	import { goto } from '$app/navigation';
	import type { DisplayError } from '../ErrorBanner.svelte';

	async function handleLogOut() {
		let logOutError: DisplayError;

		try {
			logOutError = await logOut();
		} catch (error) {
			toastStore.trigger({
				message: `Unknown Error: ${error}`,
				background: 'variant-filled-error'
			});

			return;
		}

		if (logOutError) {
			toastStore.trigger({
				message: `${logOutError?.type}: ${logOutError?.message}`,
				background: 'variant-filled-error'
			});

			return;
		}

		location.href = '/';
	}
</script>

<footer class="absolute bottom-4 left-8">
	<button type="button" class="btn variant-filled" on:click={handleLogOut}>
		<span>Log Out</span>
	</button>
</footer>
