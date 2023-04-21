<script lang="ts" context="module">
	export type DisplayError = { type: string; message: string } | null;
</script>

<script lang="ts">
	import { AlertTriangleIcon, XIcon } from 'lucide-svelte';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let error: DisplayError;
	export let dismissLabel = '';
</script>

{#if error}
	<aside class="alert variant-filled-error flex-row items-center">
		<AlertTriangleIcon class="hidden sm:block mr-4 min-w-fit" />
		<div class="alert-message">
			{#if error.type}<h3>{error.type}</h3>{/if}
			{#if error.message}<p>{error.message}</p>{/if}
		</div>
		<div class="alert-actions ml-4">
			<button
				type="button"
				class="{dismissLabel ? 'btn' : 'btn-icon'} variant-filled"
				on:click={() => dispatch('dismiss', error)}
			>
				{#if dismissLabel}{dismissLabel}{:else}<XIcon />{/if}
			</button>
		</div>
	</aside>
{/if}

<style lang="postcss">
	.alert .alert-message,
	.alert .alert-actions {
		@apply mt-0;
	}
</style>
