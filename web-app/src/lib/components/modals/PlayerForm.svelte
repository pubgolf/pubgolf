<script lang="ts" context="module">
	export type FormOperation = 'create' | 'edit';
</script>

<script lang="ts">
	import { Modal, modalStore } from '@skeletonlabs/skeleton';
	import type { PlayerData } from '$lib/proto/api/v1/shared_pb';
	import { scoringCategoryToDisplayName } from '$lib/models/scoring-category';
	import type { DisplayError } from '../ErrorBanner.svelte';
	import ErrorBanner from '../ErrorBanner.svelte';
	import type { ComponentProps } from 'svelte';
	import { XIcon } from 'lucide-svelte';

	export let parent: ComponentProps<Modal>;
	export let playerData: PlayerData;
	export let operation: FormOperation;
	export let title = '';
	export let onSubmit: (op: FormOperation, playerData: PlayerData) => Promise<DisplayError>;

	let ctaText = operation === 'create' ? 'Register Player' : 'Update Player';

	let error: DisplayError = null;
	function clearError() {
		error = null;
	}

	async function onFormSubmit() {
		if (playerData.name === '') {
			error = { type: 'Form Validation Error', message: 'Player name must not be blank.' };
			return;
		}

		const resp = await onSubmit(operation, playerData);
		if (resp) {
			error = resp;
			return;
		}

		$modalStore[0]?.response && $modalStore[0]?.response(true);
		modalStore.close();
	}
</script>

<div class="card p-4 w-modal shadow-xl space-y-4 relative">
	<header class="card-header">
		{#if title}<span class="text-2xl font-bold">{title}</span>{/if}
		<button
			type="button"
			class="btn btn-icon absolute top-4 right-4 {parent.buttonNeutral}"
			on:click={parent.onClose}><XIcon /></button
		>
	</header>

	<div class="px-4">
		<ErrorBanner {error} on:dismiss={clearError} />
	</div>

	<form class="space-y-4 p-4 pt-0">
		<label class="label">
			<span>Name</span>
			<input
				class="input"
				type="text"
				required
				placeholder="Enter name..."
				bind:value={playerData.name}
			/>
		</label>
		<label class="label">
			<span>League</span>
			<select class="select" bind:value={playerData.scoringCategory}>
				{#each Object.entries(scoringCategoryToDisplayName) as league}
					<option value={+league[0]}>{league[1]}</option>
				{/each}
			</select>
		</label>
	</form>

	<footer class="card-footer {parent.regionFooter}">
		<button class="btn {parent.buttonNeutral}" on:click={parent.onClose}
			>{parent.buttonTextCancel}</button
		>
		<button class="btn {parent.buttonPositive}" on:click={onFormSubmit}>{ctaText}</button>
	</footer>

	<slot />
</div>
