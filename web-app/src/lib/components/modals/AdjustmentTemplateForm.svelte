<script lang="ts" context="module">
	export type FormOperation = 'create' | 'edit';
</script>

<script lang="ts">
	import {
		AdjustmentTemplateDataSchema,
		type AdjustmentTemplate,
		type AdjustmentTemplateData,
		type Stage
	} from '$lib/proto/api/v1/admin_pb';
	import { create } from '@bufbuild/protobuf';
	import { Modal, modalStore } from '@skeletonlabs/skeleton';
	import { XIcon } from 'lucide-svelte';
	import type { ComponentProps } from 'svelte';
	import { writable } from 'svelte/store';
	import type { DisplayError } from '../ErrorBanner.svelte';
	import ErrorBanner from '../ErrorBanner.svelte';

	export let parent: ComponentProps<Modal>;
	export let eventKey: string;
	export let stages: Stage[];
	export let template: AdjustmentTemplate | null;
	export let operation: FormOperation;
	export let onSubmit: (data: AdjustmentTemplateData, id?: string) => Promise<DisplayError>;

	let venueSpecific = writable(!!template?.data?.stageId);
	let label = template?.data?.adjustment?.label || '';
	let amount = template?.data?.adjustment?.value || 0;
	let rank = template?.data?.rank;
	let stageId = template?.data?.stageId;
	let isVisible = template?.data?.isVisible;

	let ctaText = operation === 'create' ? 'Create Adjustment' : 'Update Adjustment';

	let error: DisplayError = null;
	function clearError() {
		error = null;
	}

	async function onFormSubmit() {
		if (label === '') {
			error = { type: 'Form Validation Error', message: 'Adjustment label must not be blank.' };
			return;
		}

		if (!amount || amount === 0) {
			error = { type: 'Form Validation Error', message: 'Adjustment value must not be zero.' };
			return;
		}

		if (!rank || rank === 0) {
			error = { type: 'Form Validation Error', message: 'Adjustment rank must not be zero.' };
			return;
		}

		const resp = await onSubmit(
			create(AdjustmentTemplateDataSchema, {
				adjustment: {
					value: amount,
					label
				},
				rank,
				eventKey,
				stageId: $venueSpecific ? stageId : undefined,
				isVisible
			}),
			template?.id
		);
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
		<span class="text-2xl font-bold"
			>{operation[0].toLocaleUpperCase() + operation.slice(1)} Adjustment</span
		>
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
			<span>Label</span>
			<input class="input" type="text" required placeholder="Enter label..." bind:value={label} />
		</label>
		<label class="label">
			<span>Value</span>
			<input
				class="input"
				type="number"
				required
				placeholder="Enter value..."
				bind:value={amount}
			/>
		</label>
		<label class="label">
			<span>Rank</span>
			<input class="input" type="number" required placeholder="Enter rank..." bind:value={rank} />
		</label>
		<label class="label">
			<input class="checkbox mr-2" type="checkbox" bind:checked={isVisible} />
			<span>Is Visible</span>
		</label>
		<label class="label">
			<input class="checkbox mr-2" type="checkbox" bind:checked={$venueSpecific} />
			<span>Venue-Specific</span>
		</label>
		{#if $venueSpecific}
			<label class="label">
				<span>Venue</span>
				<select class="select" bind:value={stageId}>
					{#each stages as stage, idx (stage.id)}
						<option value={stage.id}>{idx + 1}: {stage.venue?.name}</option>
					{/each}
				</select>
			</label>
		{/if}
	</form>

	<footer class="card-footer {parent.regionFooter}">
		<button class="btn {parent.buttonNeutral}" on:click={parent.onClose}
			>{parent.buttonTextCancel}</button
		>
		<button class="btn {parent.buttonPositive}" on:click={onFormSubmit}>{ctaText}</button>
	</footer>

	<slot />
</div>
