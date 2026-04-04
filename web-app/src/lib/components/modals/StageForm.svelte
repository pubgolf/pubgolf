<script lang="ts" context="module">
	export type FormOperation = 'create' | 'edit';
</script>

<script lang="ts">
	import type { UpdateStageRequest } from '$lib/proto/api/v1/admin_pb';
	import {
		VenueDescriptionItemType,
		ScoringCategory,
		VenueDescriptionItemSchema
	} from '$lib/proto/api/v1/shared_pb';
	import type { Venue } from '$lib/proto/api/v1/shared_pb';
	import { Modal, modalStore } from '@skeletonlabs/skeleton';
	import { ArrowUpIcon, ArrowDownIcon, TrashIcon, PlusIcon, XIcon } from 'lucide-svelte';
	import type { ComponentProps } from 'svelte';
	import type { DisplayError } from '../ErrorBanner.svelte';
	import ErrorBanner from '../ErrorBanner.svelte';
	import { create } from '@bufbuild/protobuf';

	export let parent: ComponentProps<Modal>;
	export let venues: Venue[];
	export let stage: UpdateStageRequest;
	export let onSubmit: (data: UpdateStageRequest) => Promise<DisplayError>;

	// Local mutable copy of items for form state.
	let items: {
		content: string;
		itemType: VenueDescriptionItemType;
		audiences: ScoringCategory[];
	}[] = (stage.venueDescriptions ?? []).map((item) => ({
		content: item.content,
		itemType: item.itemType,
		audiences: [...item.audiences]
	}));

	const itemTypeLabels: Record<number, string> = {
		[VenueDescriptionItemType.DEFAULT]: 'Default',
		[VenueDescriptionItemType.WARNING]: 'Warning',
		[VenueDescriptionItemType.RULE]: 'Rule'
	};

	const audienceLabels: { value: ScoringCategory; label: string }[] = [
		{ value: ScoringCategory.PUB_GOLF_NINE_HOLE, label: '9-Hole' },
		{ value: ScoringCategory.PUB_GOLF_FIVE_HOLE, label: '5-Hole' },
		{ value: ScoringCategory.PUB_GOLF_CHALLENGES, label: 'Challenges' }
	];

	let error: DisplayError = null;
	function clearError() {
		error = null;
	}

	function addItem() {
		items = [...items, { content: '', itemType: VenueDescriptionItemType.DEFAULT, audiences: [] }];
	}

	function removeItem(idx: number) {
		items.splice(idx, 1);
		items = items;
	}

	function moveUp(idx: number) {
		if (idx <= 0) return;
		[items[idx - 1], items[idx]] = [items[idx], items[idx - 1]];
		items = items;
	}

	function moveDown(idx: number) {
		if (idx >= items.length - 1) return;
		[items[idx], items[idx + 1]] = [items[idx + 1], items[idx]];
		items = items;
	}

	function toggleAudience(itemIdx: number, category: ScoringCategory) {
		const auds = items[itemIdx].audiences;
		const existingIdx = auds.indexOf(category);

		if (existingIdx >= 0) {
			auds.splice(existingIdx, 1);
		} else {
			auds.push(category);
		}

		items = items;
	}

	function onItemTypeChange(idx: number) {
		if (items[idx].itemType !== VenueDescriptionItemType.RULE) {
			items[idx].audiences = [];
			items = items;
		}
	}

	async function onFormSubmit() {
		if (!stage.durationMin || stage.durationMin === 0) {
			error = { type: 'Form Validation Error', message: 'Duration value must not be zero.' };
			return;
		}

		stage.venueDescriptions = items.map((item) =>
			create(VenueDescriptionItemSchema, {
				content: item.content,
				itemType: item.itemType,
				audiences: item.audiences
			})
		);
		stage.venueDescription = '';

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

		<div class="space-y-2">
			<span class="label"><span>Rule Items</span></span>

			{#each items as item, idx (idx)}
				<div class="card p-3 variant-soft space-y-2">
					<div class="flex items-center gap-2">
						<select
							class="select flex-1"
							bind:value={item.itemType}
							on:change={() => onItemTypeChange(idx)}
						>
							{#each Object.entries(itemTypeLabels) as [value, label]}
								<option value={Number(value)}>{label}</option>
							{/each}
						</select>

						<button
							type="button"
							class="btn btn-icon btn-sm variant-ghost"
							disabled={idx === 0}
							on:click={() => moveUp(idx)}
						>
							<ArrowUpIcon size={16} />
						</button>
						<button
							type="button"
							class="btn btn-icon btn-sm variant-ghost"
							disabled={idx === items.length - 1}
							on:click={() => moveDown(idx)}
						>
							<ArrowDownIcon size={16} />
						</button>
						<button
							type="button"
							class="btn btn-icon btn-sm variant-ghost-error"
							on:click={() => removeItem(idx)}
						>
							<TrashIcon size={16} />
						</button>
					</div>

					<textarea
						class="textarea min-h-[80px]"
						placeholder="Enter item text..."
						bind:value={item.content}
					></textarea>

					{#if item.itemType === VenueDescriptionItemType.RULE}
						<div class="flex gap-3 items-center flex-wrap">
							<span class="text-sm font-medium">Audiences:</span>
							{#each audienceLabels as { value, label }}
								<label class="flex items-center gap-1">
									<input
										type="checkbox"
										class="checkbox"
										checked={item.audiences.includes(value)}
										on:change={() => toggleAudience(idx, value)}
									/>
									<span class="text-sm">{label}</span>
								</label>
							{/each}
						</div>
					{/if}
				</div>
			{/each}

			<button type="button" class="btn variant-ghost-primary w-full" on:click={addItem}>
				<PlusIcon size={16} />
				<span>Add Item</span>
			</button>
		</div>
	</form>

	<footer class="card-footer {parent.regionFooter}">
		<button class="btn {parent.buttonNeutral}" on:click={parent.onClose}
			>{parent.buttonTextCancel}</button
		>
		<button class="btn {parent.buttonPositive}" on:click={onFormSubmit}>Update Stage</button>
	</footer>

	<slot />
</div>
