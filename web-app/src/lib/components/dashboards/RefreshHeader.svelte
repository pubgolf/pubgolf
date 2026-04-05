<script lang="ts">
	import { formatTimestamp } from '$lib/helpers/formatters';
	import type { Snippet } from 'svelte';

	import { RefreshCwIcon } from 'lucide-svelte';

	interface Props {
		title: string;
		loadingStatus: Promise<void>;
		lastRefresh?: Date;
		refresh: (() => void) | (() => Promise<void>);
		filters?: Snippet;
	}

	let { title, loadingStatus, lastRefresh, refresh, filters }: Props = $props();

	let state: 'ready' | 'loading' | 'error' = $state('loading');
	$effect(() => {
		let cancelled = false;
		state = 'loading';
		loadingStatus
			.then(() => !cancelled && (state = 'ready'))
			.catch(() => !cancelled && (state = 'error'));
		return () => {
			cancelled = true;
		};
	});
</script>

<header class="flex justify-between items-start mb-4 md:mt-4">
	<div>
		<h1 class="mb-2">{title}</h1>
		{@render filters?.()}
	</div>
	<div class="text-right">
		<button type="button" class="btn variant-filled mb-0.5" onclick={refresh}>
			<span class="badge-icon"
				><RefreshCwIcon class={state === 'loading' ? 'animate-spin' : ''} /></span
			>
			<span>Refresh</span>
		</button><br />
		{#if state === 'loading'}
			<span class="text-xs">Fetching data...</span>
		{:else if state === 'ready'}
			{#if lastRefresh}
				<span class="text-xs">Last Fetched: {formatTimestamp(lastRefresh)}</span>
			{/if}
		{:else}
			<span class="text-xs">Error fetching data</span>
		{/if}
	</div>
</header>
