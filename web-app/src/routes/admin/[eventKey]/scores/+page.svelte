<script lang="ts">
	import { page } from '$app/stores';
	import ErrorBanner from '$lib/components/ErrorBanner.svelte';
	import NoDataCard from '$lib/components/dashboards/NoDataCard.svelte';
	import RefreshHeader from '$lib/components/dashboards/RefreshHeader.svelte';
	import NewScoreForm from '$lib/components/modals/ScoreForm.svelte';
	import SetTitle from '$lib/components/util/SetTitle.svelte';
	import { formatPlayerName as playerNameWithLeague } from '$lib/helpers/formatters';
	import { combineIds, separateIds } from '$lib/helpers/scores';
	import type { Stage, StageScore } from '$lib/proto/api/v1/admin_pb';
	import type { Player } from '$lib/proto/api/v1/shared_pb';
	import { getAdminServiceClient } from '$lib/rpc/client';
	import { modalStore, toastStore } from '@skeletonlabs/skeleton';
	import { BeerIcon, RefreshCwIcon } from 'lucide-svelte';
	import type { ComponentProps } from 'svelte';
	import { onMount } from 'svelte';
	import { noop } from 'svelte/internal';

	let refreshProgress: Promise<void> = new Promise(noop);
	let lastRefresh: Date;

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
		refreshProgress = Promise.all([fetchPlayers(), fetchStages(), fetchScores()]).then(noop);
		await refreshProgress;
		lastRefresh = new Date();
	}
	onMount(refreshData);

	function getVenueName(stageId: string) {
		const i = stages.findIndex((x) => x.id === stageId);
		return `${i + 1}: ${stages[i].venue?.name}`;
	}

	function getPlayerName(id: string) {
		const player = players.find((x) => x.id === id);
		if (!player) {
			return '[Unknown Player]';
		}

		return playerNameWithLeague(player, $page.params.eventKey);
	}

	async function deleteScore(score: StageScore) {
		const rpc = await getAdminServiceClient();
		try {
			await rpc.deleteStageScore({
				stageId: score.stageId,
				playerId: score.playerId
			});
		} catch (error) {
			toastStore.trigger({
				message: `API Error: ${error}`,
				background: 'variant-filled-error'
			});
		}

		refreshData();
	}

	async function attemptDeleteScore(score: StageScore) {
		modalStore.trigger({
			type: 'confirm',
			title: 'Confirm Deletion',
			body: `Are you sure you wish to delete the score at ${getVenueName(
				score.stageId
			)} for ${getPlayerName(score.playerId)}?`,
			response: (r: boolean) => r && deleteScore(score)
		});
	}

	async function showNewScoreModal() {
		const props: Omit<ComponentProps<NewScoreForm>, 'parent'> = {
			eventKey: $page.params.eventKey,
			players: await players,
			stages: await stages,
			onSubmit: async (data) => {
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
			eventKey: $page.params.eventKey,
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
</script>

<SetTitle title="Scores" />

<div class="max-w-3xl mx-auto mb-4">
	<RefreshHeader
		title="Scores"
		refresh={refreshData}
		loadingStatus={refreshProgress}
		{lastRefresh}
	/>

	{#await refreshProgress}
		<div class="card py-12 flex flex-col items-center">
			<p class="mb-4">Loading player scores...</p>
			<RefreshCwIcon class="animate-spin" />
		</div>
	{:then}
		{#if scores.length}
			<div class="table-container">
				<table class="table table-hover">
					<thead>
						<tr>
							<th>Venue</th>
							<th>Player</th>
							<th>Score</th>
							<th class="table-cell-fit">Actions</th>
						</tr>
					</thead>
					<tbody>
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
								<td class="table-cell-fit action-btns">
									<button
										type="button"
										class="btn btn-sm variant-filled"
										on:click={() => showEditScoreModal(score)}
									>
										<span>Edit</span>
									</button>
									<button
										type="button"
										class="btn btn-sm variant-filled-error ml-4"
										on:click={() => attemptDeleteScore(score)}
									>
										<span>Delete</span>
									</button>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{:else}
			<NoDataCard
				text="No Scores to Display"
				ctaText="Enter a Score"
				on:click={showNewScoreModal}
			/>
		{/if}
	{:catch error}
		<ErrorBanner
			error={{ type: 'Server Error', message: error }}
			dismissLabel="Retry"
			on:dismiss={refreshData}
		/>
	{/await}

	<footer class="fixed bottom-8 right-4">
		<button type="button" class="btn btn-lg variant-filled-primary" on:click={showNewScoreModal}>
			<span class="badge-icon"><BeerIcon /></span>
			<span>Enter Score</span>
		</button>
	</footer>
</div>

<style lang="postcss">
	.table tbody td {
		@apply align-middle;
	}
	.table tbody td.action-btns {
		@apply whitespace-nowrap;
	}
</style>
