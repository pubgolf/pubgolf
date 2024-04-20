<script lang="ts" context="module">
	export type FormOperation = 'create' | 'edit';
</script>

<script lang="ts">
	import type { UpdateStageRequest } from '$lib/proto/api/v1/admin_pb';
	import type { Venue } from '$lib/proto/api/v1/shared_pb';
	import { Modal, modalStore } from '@skeletonlabs/skeleton';
	import { XIcon } from 'lucide-svelte';
	import type { ComponentProps } from 'svelte';
	import type { DisplayError } from '../ErrorBanner.svelte';
	import ErrorBanner from '../ErrorBanner.svelte';

	export let parent: ComponentProps<Modal>;
	export let venues: Venue[];
	export let stage: UpdateStageRequest;
	export let onSubmit: (data: UpdateStageRequest) => Promise<DisplayError>;

	let error: DisplayError = null;
	function clearError() {
		error = null;
	}

	async function onFormSubmit() {
		if (!stage.durationMin || stage.durationMin === 0) {
			error = { type: 'Form Validation Error', message: 'Duration value must not be zero.' };
			return;
		}

		const resp = await onSubmit(stage);
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
		<span class="text-2xl font-bold">Update Stage</span>
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
			<span>Venue</span>
			<select class="select" bind:value={stage.venueId}>
				{#each venues as venue (venue.id)}
					<option value={venue.id}>{venue.name}</option>
				{/each}
			</select>
		</label>
		<label class="label">
			<span>Duration (min)</span>
			<input
				class="input"
				type="number"
				required
				placeholder="Enter duration..."
				bind:value={stage.durationMin}
			/>
		</label>
		<label class="label">
			<span>Rule</span>
			<textarea
				class="textarea min-h-[200px]"
				required
				placeholder="Enter rule text..."
				bind:value={stage.venueDescription}
			/>
		</label>
	</form>

	<footer class="card-footer {parent.regionFooter}">
		<button class="btn {parent.buttonNeutral}" on:click={parent.onClose}
			>{parent.buttonTextCancel}</button
		>
		<button class="btn {parent.buttonPositive}" on:click={onFormSubmit}>Update Stage</button>
	</footer>

	<slot />
</div>
