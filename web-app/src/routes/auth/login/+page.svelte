<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { logIn } from '$lib/auth/client';
	import type { DisplayError } from '$lib/components/ErrorBanner.svelte';
	import ErrorBanner from '$lib/components/ErrorBanner.svelte';
	import { onMount } from 'svelte';

	let passwordInput: HTMLInputElement;
	onMount(() => passwordInput.focus());

	let password: string;
	let error: DisplayError = null;
	function clearError() {
		error = null;
	}

	async function onFormSubmit() {
		let e = await logIn(password);
		if (e) {
			error = e;
			return;
		}

		const redirectURL = new URL($page.url.searchParams.get('redirect') || '/', $page.url);
		if (redirectURL.origin === $page.url.origin) {
			goto(redirectURL.pathname);
		} else {
			goto('/');
		}
	}
</script>

<div class="flex items-start md:items-center md:justify-center h-screen">
	<div class="card w-full max-w-xl m-4">
		<header class="card-header mb-4">
			<span class="text-2xl font-bold">Admin Log In</span>
		</header>

		<div class="m-4 mt-0" class:hidden={!error}>
			<ErrorBanner {error} on:dismiss={clearError} />
		</div>

		<section class="p-4 pt-0">
			<form class="space-y-4 mb-4" on:submit={onFormSubmit}>
				<label class="label">
					<span>Username</span>
					<input class="input" type="text" readonly value="pubgolf_admin" />
				</label>
				<label class="label">
					<span>Password</span>
					<input
						bind:this={passwordInput}
						class="input"
						type="password"
						required
						placeholder="Enter name..."
						bind:value={password}
					/>
				</label>
			</form>
		</section>
		<footer class="card-footer flex justify-end">
			<button class="btn variant-filled" on:click={onFormSubmit}>Log In</button>
		</footer>
	</div>
</div>
