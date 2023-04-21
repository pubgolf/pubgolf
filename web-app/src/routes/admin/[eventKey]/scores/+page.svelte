<script lang="ts">
	import { BeerIcon, RefreshCcwIcon } from 'lucide-svelte';
	import SetTitle from '$lib/components/util/SetTitle.svelte';
	import NewScoreForm from '$lib/components/modals/NewScoreForm.svelte';
	import { modalStore } from '@skeletonlabs/skeleton';
	import type { ComponentProps } from 'svelte';
	import { getAdminServiceClient } from '$lib/rpc/client';
	import type { Player } from '$lib/proto/api/v1/shared_pb';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import type { Stage, StageScore } from '$lib/proto/api/v1/admin_pb';
	import { scoringCategoryToDisplayName } from '$lib/models/scoring-category';
	import ErrorBanner from '$lib/components/ErrorBanner.svelte';
	import { combineIds, separateIds } from '$lib/models/scores';

	let dataReady: Promise<any> = new Promise(() => {});
	let dataUpdatedAt: Date;
	function formatTimestamp(d: Date) {
		return new Intl.DateTimeFormat(undefined, { timeStyle: 'medium' }).format(d);
	}

	let players: Player[] = [];
	async function fetchPlayers() {
		const rpc = await getAdminServiceClient();
		const resp = await rpc.listPlayers({
			eventKey: $page.params.eventKey
		});

		players = resp.players;
	}

	let stages: Stage[] = [];
	async function fetchStages() {
		const rpc = await getAdminServiceClient();
		const resp = await rpc.listEventStages({
			eventKey: $page.params.eventKey
		});

		stages = resp.stages;
	}

	let scores: StageScore[] = [];
	async function fetchScores() {
		const rpc = await getAdminServiceClient();
		const resp = await rpc.listStageScores({
			eventKey: $page.params.eventKey
		});

		scores = resp.scores.reverse();
	}

	async function refreshData() {
		dataReady = Promise.all([fetchPlayers(), fetchStages(), fetchScores()]);
		await dataReady;
		dataUpdatedAt = new Date();
	}

	function getVenueName(stageId: string) {
		const i = stages.findIndex((x) => x.id === stageId);
		return `Venue ${i + 1}: ${stages[i].venue?.name}`;
	}

	function getPlayerName(id: string) {
		return players.find((x) => x.id === id)?.data?.name;
	}

	async function showNewScoreModal() {
		const props: Omit<ComponentProps<NewScoreForm>, 'parent'> = {
			players: await players,
			stages: await stages,
			onSubmit: async (data, _) => {
				const rpc = await getAdminServiceClient();
				try {
					await rpc.createStageScore({
						data
					});
				} catch (error) {
					return { type: 'Server Error', message: error as string };
				}

				return null;
			}
		};

		modalStore.trigger({
			type: 'component',
			component: {
				ref: NewScoreForm,
				props
			},
			response: (submittedForm: boolean) => submittedForm && refreshData()
		});
	}

	async function showEditScoreModal(score: StageScore) {
		const props: Omit<ComponentProps<NewScoreForm>, 'parent'> = {
			title: 'Edit Score',
			ctaText: 'Save',
			players: await players,
			stages: await stages,
			score: separateIds(score),
			onSubmit: async (data, ids) => {
				if (!ids) {
					return {
						type: 'Application Error',
						message: 'Expected to receive score IDs to submit with edit request.'
					};
				}

				const rpc = await getAdminServiceClient();
				try {
					await rpc.updateStageScore({ score: combineIds({ data, ids }) });
				} catch (error) {
					return { type: 'Server Error', message: error as string };
				}

				return null;
			}
		};

		modalStore.trigger({
			type: 'component',
			component: {
				ref: NewScoreForm,
				props
			},
			response: (submittedForm: boolean) => submittedForm && refreshData()
		});
	}

	onMount(refreshData);
</script>

<SetTitle title="Scores" />

<div class="max-w-3xl mx-auto">
	<header class="flex justify-between items-start mb-4 md:mt-4">
		<h1>Scores</h1>
		<div class="text-right">
			<button type="button" class="btn variant-filled mb-0.5" on:click={refreshData}>
				<span class="badge-icon"><RefreshCcwIcon /></span>
				<span>Refresh</span>
			</button><br />
			{#await dataReady}
				<span class="text-xs">Fetching data...</span>
			{:then}
				<span class="text-xs">Last Fetched: {formatTimestamp(dataUpdatedAt)}</span>
			{:catch error}
				<span class="text-xs">Error fetching data</span>
			{/await}
		</div>
	</header>

	<div class="table-container">
		{#await dataReady}
			<div class="py-12">
				<p class="mb-4">Fetching data...</p>
			</div>
		{:then}
			<table class="table table-hover">
				<thead>
					<tr>
						<th>Venue</th>
						<th>Player</th>
						<th>Score</th>
						<th class="table-cell-fit">Edit</th>
					</tr>
				</thead>
				<tbody>
					{#if scores.length}
						{#each scores as score (score.score?.id)}
							<tr>
								<td>{getVenueName(score.stageId)}</td>
								<td>{getPlayerName(score.playerId)}</td>
								<td class="font-mono">
									Sips: {score.score?.data?.value}
									{#if score.adjustments.length > 0}
										{#each score.adjustments as adj (adj.id)}
											<br /><span class="pl-4"
												>{adj.data?.value && adj.data?.value > 0 ? '+' : ''}{adj.data?.value}: {adj
													.data?.label}</span
											>
										{/each}
									{/if}
								</td>
								<td class="table-cell-fit">
									<button
										type="button"
										class="btn btn-sm variant-filled"
										on:click={() => showEditScoreModal(score)}
									>
										<span>Edit</span>
									</button>
								</td>
							</tr>
						{/each}
					{:else}
						<tr class="text-center">
							<td colspan="3">
								<div class="py-12">
									<p class="mb-4">No Scores to Display</p>
									<button
										type="button"
										class="btn btn-lg variant-filled-secondary"
										on:click={showNewScoreModal}
									>
										<span>Enter a Score</span>
									</button>
								</div>
							</td>
						</tr>
					{/if}
				</tbody>
			</table>
		{:catch error}
			<ErrorBanner
				error={{ type: 'Server Error', message: error }}
				dismissLabel="Retry"
				on:dismiss={refreshData}
			/>
		{/await}
	</div>

	<footer class="fixed bottom-8 right-4">
		<button type="button" class="btn btn-lg variant-filled-primary" on:click={showNewScoreModal}>
			<span class="badge-icon"><BeerIcon /></span>
			<span>Enter Score</span>
		</button>
	</footer>
</div>
