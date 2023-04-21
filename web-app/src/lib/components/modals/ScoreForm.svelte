<script lang="ts">
	import { formatListAnd } from '$lib/helpers/formatters';
	import type { StageScoreIds, Strict } from '$lib/helpers/scores';
	import { scoringCategoryToDisplayName } from '$lib/helpers/scoring-category';
	import type { AdjustmentData, Stage, StageScoreData } from '$lib/proto/api/v1/admin_pb';
	import type { Player } from '$lib/proto/api/v1/shared_pb';
	import type { PlainMessage } from '@bufbuild/protobuf';
	import { modalStore, type Modal } from '@skeletonlabs/skeleton';
	import { PlusIcon, XIcon } from 'lucide-svelte';
	import type { ComponentProps } from 'svelte';
	import ErrorBanner, { type DisplayError } from '../ErrorBanner.svelte';

	function selectOnFocus(e: FocusEvent) {
		(e.target as HTMLInputElement | null)?.select();
	}

	export let parent: ComponentProps<Modal>;
	export let title = 'Enter a Score';
	export let ctaText = 'Create Score';
	export let players: Player[];
	export let stages: Stage[];
	export let score: {
		data: Strict<StageScoreData>;
		ids: StageScoreIds;
	} | null = null;
	export let onSubmit: (
		scoreData: Strict<StageScoreData>,
		ids?: StageScoreIds
	) => Promise<DisplayError>;

	let editMode = !!score;

	let formInit = false;
	$: {
		if (!formInit && score) {
			formData = score.data;

			let newBonuses: AdjustmentFormEntry[] = [];
			let newPenalties: AdjustmentFormEntry[] = [];
			score.data.adjustments.forEach((a, i) => {
				if (a.value < 0) {
					newBonuses.push({
						idType: 'server',
						id: score?.ids.adjustments[i].id || '',
						adjustment: {
							label: a.label,
							value: -1 * a.value
						}
					});
				}

				if (a.value > 0) {
					newPenalties.push({
						idType: 'server',
						id: score?.ids.adjustments[i].id || '',
						adjustment: a
					});
				}
			});

			bonuses = newBonuses;
			penalties = newPenalties;

			formInit = true;
		}
	}

	type AdjustmentFormEntry = {
		idType: 'server' | 'client';
		id: string;
		adjustment: PlainMessage<AdjustmentData>;
	};

	let penalties: AdjustmentFormEntry[] = [];
	async function addPenalty() {
		penalties = [
			...penalties,
			{
				idType: 'client',
				id: crypto.randomUUID(),
				adjustment: {
					value: 0,
					label: ''
				}
			}
		];
	}
	function removePenalty(id: string) {
		return () => {
			const idx = penalties.findIndex((x) => x.id === id);
			if (idx > -1) {
				penalties.splice(idx, 1);
				penalties = penalties;
			}
		};
	}

	let bonuses: AdjustmentFormEntry[] = [];
	async function addBonus() {
		bonuses = [
			...bonuses,
			{
				idType: 'client',
				id: crypto.randomUUID(),
				adjustment: {
					value: 0,
					label: ''
				}
			}
		];
	}
	function removeBonus(id: string) {
		return () => {
			const idx = bonuses.findIndex((x) => x.id === id);
			if (idx > -1) {
				bonuses.splice(idx, 1);
				bonuses = bonuses;
			}
		};
	}

	$: formData.adjustments = [
		...penalties.map((x) => x.adjustment),
		...bonuses.map((x) => ({ ...x.adjustment, value: -x.adjustment.value }))
	];

	let formData: Strict<StageScoreData> = {
		stageId: '',
		playerId: '',
		score: {
			value: 0
		},
		adjustments: []
	};

	let error: DisplayError = null;
	function clearError() {
		error = null;
	}

	async function onFormSubmit() {
		let blankFields: string[] = [];

		if (formData.playerId === '') {
			blankFields.push('Player');
		}

		if (formData.stageId === '') {
			blankFields.push('Venue');
		}

		formData.adjustments.forEach((x, i) => {
			if (x.label === '') {
				if (i < penalties.length) {
					blankFields.push(`Label for penalty #${i + 1}`);
				} else {
					blankFields.push(`Label for bonus #${i + 1 - penalties.length}`);
				}
			}
		});

		if (blankFields.length) {
			error = {
				type: 'Form Error',
				message: `The field${blankFields.length > 1 ? 's' : ''} ${formatListAnd(
					blankFields.map((x) => `"${x}"`)
				)} must not be blank.`
			};
			return;
		}

		let zeroFields: string[] = [];
		if (formData.score.value === 0) {
			zeroFields.push('Score');
		}

		formData.adjustments.forEach((x, i) => {
			if (x.value === 0) {
				if (i < penalties.length) {
					zeroFields.push(`${x.label} Amount`);
				} else {
					zeroFields.push(`${x.label} Amount`);
				}
			}
		});

		if (zeroFields.length) {
			let fields = zeroFields.map((x) => `"${x}"`);
			let message = `${fields[0]} must be a non-zero value.`;
			if (fields.length > 1) {
				message = `The fields ${formatListAnd(fields)} must be non-zero values.`;
			}

			error = {
				type: 'Form Error',
				message
			};
			return;
		}

		let ids: StageScoreIds | undefined;
		if (score?.ids) {
			ids = {
				score: {
					id: score.ids.score.id
				},
				adjustments: [
					...penalties.map((x) => {
						return {
							id: x.idType === 'server' ? x.id : ''
						};
					}),
					...bonuses.map((x) => {
						return {
							id: x.idType === 'server' ? x.id : ''
						};
					})
				]
			};
		}

		const resp = await onSubmit(formData, ids);
		if (resp) {
			error = resp;
			return;
		}

		$modalStore[0]?.response && $modalStore[0]?.response(true);
		modalStore.close();
	}
</script>

<div class="card p-4 w-modal shadow-xl space-y-4 relative">
	{#if title}
		<header class="card-header">
			<span class="text-2xl font-bold">{title}</span>
			<button
				type="button"
				class="btn btn-icon absolute top-4 right-4 {parent.buttonNeutral}"
				on:click={parent.onClose}><XIcon /></button
			>
		</header>
	{/if}

	<form class="space-y-4 p-2 sm:p-4 pt-0">
		<div class="grid sm:grid-cols-2 gap-4 mb-8">
			<label class="label">
				<span>Player</span>
				<select class="select" required disabled={editMode} bind:value={formData.playerId}>
					{#each players as player}
						<option value={player.id}
							>{player.data?.name} ({player.data?.scoringCategory
								? scoringCategoryToDisplayName[player.data?.scoringCategory]
								: 'Not Set'})</option
						>
					{/each}
				</select>
			</label>
			<label class="label">
				<span>Venue</span>
				<select class="select" disabled={editMode} bind:value={formData.stageId}>
					{#each stages as stage, idx (stage.id)}
						<option value={stage.id}>{stage.venue?.name || `Venue #${idx + 1}`}</option>
					{/each}
				</select>
			</label>
			<label class="label">
				<span>Score</span>
				<input
					class="input"
					type="number"
					inputmode="numeric"
					required
					on:focus={selectOnFocus}
					bind:value={formData.score.value}
				/>
			</label>
		</div>

		<span class="h3 block">Penalties</span>
		{#each penalties as penalty, idx (penalty.id)}
			<div class="flex items-end gap-2">
				<button
					type="button"
					class="btn btn-icon variant-ringed-secondary shrink-0"
					on:click={removePenalty(penalty.id)}><XIcon /></button
				>
				<label class="label grow">
					<span class:sr-only={idx > 0}>Label</span>{#if idx == 0}<br />{/if}
					<!-- svelte-ignore a11y-autofocus -->
					<input
						class="input"
						type="text"
						placeholder="Reason"
						required
						bind:value={penalty.adjustment.label}
						autofocus={idx === penalties.length - 1}
					/>
				</label>
				<label class="label w-24">
					<span class:sr-only={idx > 0}>Amount</span>
					<input
						class="input"
						type="number"
						inputmode="numeric"
						required
						on:focus={selectOnFocus}
						bind:value={penalty.adjustment.value}
					/>
				</label>
			</div>
		{/each}
		<button type="button" class="btn {parent.buttonNeutral}" on:click={addPenalty}>
			<span class="badge-icon mr-2"><PlusIcon /></span>
			Add Penalty</button
		>

		<span class="h3 block">Bonuses</span>
		{#each bonuses as bonus, idx (bonus.id)}
			<div class="flex items-end gap-2">
				<button
					type="button"
					class="btn btn-icon variant-ringed-secondary shrink-0"
					on:click={removeBonus(bonus.id)}><XIcon /></button
				>
				<label class="label grow">
					<span class:sr-only={idx > 0}>Label</span>{#if idx == 0}<br />{/if}
					<!-- svelte-ignore a11y-autofocus -->
					<input
						class="input"
						type="text"
						placeholder="Reason"
						required
						bind:value={bonus.adjustment.label}
						autofocus={idx === bonuses.length - 1}
					/>
				</label>
				<label class="label w-24">
					<span class:sr-only={idx > 0}>Amount</span>
					<input
						class="input"
						type="number"
						inputmode="numeric"
						required
						on:focus={selectOnFocus}
						bind:value={bonus.adjustment.value}
					/>
				</label>
			</div>
		{/each}
		<button type="button" class="btn {parent.buttonNeutral}" on:click={addBonus}>
			<span class="badge-icon mr-2"><PlusIcon /></span>
			Add Bonus</button
		>
	</form>

	<div class="px-4">
		<ErrorBanner {error} on:dismiss={clearError} />
	</div>

	<footer class="card-footer {parent.regionFooter} pt-4">
		<div>
			<button class="btn {parent.buttonNeutral}" on:click={parent.onClose}
				>{parent.buttonTextCancel}</button
			>
			<button class="btn {parent.buttonPositive}" on:click={onFormSubmit}>{ctaText}</button>
		</div>
	</footer>

	<slot />
</div>
