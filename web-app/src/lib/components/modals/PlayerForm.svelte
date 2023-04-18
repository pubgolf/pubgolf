<script lang="ts" context="module">
	export type FormOperation = 'create' | 'edit';
	export type FormError = { type: string; message: string } | null;
</script>

<script lang="ts">
	import { modalStore } from '@skeletonlabs/skeleton';
	import { AlertTriangleIcon, XIcon } from 'lucide-svelte';
	import type { PlayerData } from '$lib/proto/api/v1/shared_pb';
	import { scoringCategoryToDisplayName } from '$lib/models/scoring-category';

	export let parent: any;
	export let playerData: PlayerData;
	export let operation: FormOperation;
	export let onSubmit: (op: FormOperation, playerData: PlayerData) => Promise<FormError>;

	let error: FormError = null;
	let ctaText = operation === 'create' ? 'Register Player' : 'Update Player';

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

<div class="modal-example-form card p-4 w-modal shadow-xl space-y-4 relative">
	{#if $modalStore[0]?.title}
		<header class="text-2xl font-bold">{$modalStore[0]?.title}</header>
	{/if}

	{#if error}
		<aside class="alert variant-filled-error flex-row items-center">
			<AlertTriangleIcon class="hidden sm:block mr-4" />
			<div class="alert-message">
				<h3>{error.type}</h3>
				<p>{error.message}</p>
			</div>
			<div class="alert-actions">
				<button
					type="button"
					class="btn-icon variant-filled"
					on:click={() => {
						error = null;
					}}><XIcon /></button
				>
			</div>
		</aside>
	{/if}

	{#if $modalStore[0]?.body}
		<article>{$modalStore[0]?.body}</article>
	{/if}

	<form class="modal-form space-y-4 mb-8">
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

	<footer class="modal-footer {parent.regionFooter}">
		<button class="btn {parent.buttonNeutral}" on:click={parent.onClose}
			>{parent.buttonTextCancel}</button
		>
		<button class="btn {parent.buttonPositive}" on:click={onFormSubmit}>{ctaText}</button>
	</footer>
</div>

<style lang="postcss">
	.alert .alert-message,
	.alert .alert-actions {
		@apply mt-0;
	}
</style>
