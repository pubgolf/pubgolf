<script lang="ts" context="module">
	export type FormOperation = 'create' | 'edit';
</script>

<script lang="ts">
	import { scoringCategoryToDisplayName } from '$lib/helpers/scoring-category';
	import { PlayerData, ScoringCategory, type Player } from '$lib/proto/api/v1/shared_pb';
	import { Modal, modalStore } from '@skeletonlabs/skeleton';
	import { XIcon } from 'lucide-svelte';
	import type { ComponentProps } from 'svelte';
	import type { DisplayError } from '../ErrorBanner.svelte';
	import ErrorBanner from '../ErrorBanner.svelte';

	export let parent: ComponentProps<Modal>;
	export let eventKey: string;
	export let player: Player;
	export let operation: FormOperation;
	export let title = '';
	export let onSubmit: (
		op: FormOperation,
		playerData: PlayerData,
		scoringCategory: ScoringCategory,
		phoneNumber?: string
	) => Promise<DisplayError>;

	let playerName = player.data?.name || '';
	let phoneNumber = '';
	let scoringCategory =
		player.events.find((x) => x.eventKey == eventKey)?.scoringCategory ||
		ScoringCategory.UNSPECIFIED;

	let ctaText = operation === 'create' ? 'Register Player' : 'Update Player';
	let requirePhoneNumber = operation === 'create';
	let normalizedPhoneNumber = '';
	$: normalizedPhoneNumber = cleanPhoneNumber(phoneNumber);

	let error: DisplayError = null;
	function clearError() {
		error = null;
	}

	function cleanPhoneNumber(num: string): string {
		num = num.replaceAll(/[^\d]/g, '');
		if (num.length < 11 && !num.startsWith('1')) {
			num = '1' + num;
		}

		return '+' + num;
	}

	async function onFormSubmit() {
		if (playerName === '') {
			error = { type: 'Form Validation Error', message: 'Player name must not be blank.' };
			return;
		}

		if (requirePhoneNumber && normalizedPhoneNumber.length !== 12) {
			error = { type: 'Form Validation Error', message: 'Phone number must be 10 digits.' };
			return;
		}

		const resp = await onSubmit(
			operation,
			new PlayerData({
				name: playerName
			}),
			scoringCategory,
			requirePhoneNumber ? normalizedPhoneNumber : undefined
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
				bind:value={playerName}
			/>
		</label>
		{#if requirePhoneNumber}
			<label class="label">
				<span>Phone Number</span>
				<input class="input" type="tel" required bind:value={phoneNumber} />
			</label>
		{/if}
		<label class="label">
			<span>League</span>
			<select class="select" bind:value={scoringCategory}>
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
