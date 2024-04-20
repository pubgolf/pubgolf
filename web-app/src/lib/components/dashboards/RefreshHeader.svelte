<script lang="ts">
	import { formatTimestamp } from '$lib/helpers/formatters';

	import { RefreshCwIcon } from 'lucide-svelte';

	export let title: string;
	export let loadingStatus: Promise<void>;
	export let lastRefresh: Date | undefined = undefined;
	export let refresh: (() => void) | (() => Promise<void>);

	let state: 'ready' | 'loading' | 'error' = 'loading';
	function setState(loadingStatus: Promise<void>) {
		if (loadingStatus) {
			state = 'loading';
			loadingStatus.then(() => (state = 'ready')).catch(() => (state = 'error'));
		}
	}
	$: setState(loadingStatus);
</script>

<header class="flex justify-between items-start mb-4 md:mt-4">
	<div>
		<h1 class="mb-2">{title}</h1>
		<slot name="filters" />
	</div>
	<div class="text-right">
		<button type="button" class="btn variant-filled mb-0.5" on:click={refresh}>
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
