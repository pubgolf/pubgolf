<script lang="ts">
	import { page } from '$app/stores';
	import type { DisplayError } from '$lib/components/ErrorBanner.svelte';
	import ErrorBanner from '$lib/components/ErrorBanner.svelte';
	import NoDataCard from '$lib/components/dashboards/NoDataCard.svelte';
	import RefreshHeader from '$lib/components/dashboards/RefreshHeader.svelte';
	import PlayerForm, { type FormOperation } from '$lib/components/modals/PlayerForm.svelte';
	import SetTitle from '$lib/components/util/SetTitle.svelte';
	import { scoringCategoryToDisplayName } from '$lib/helpers/scoring-category';
	import { PlayerData, type Player } from '$lib/proto/api/v1/shared_pb';
	import { getAdminServiceClient } from '$lib/rpc/client';
	import { modalStore } from '@skeletonlabs/skeleton';
	import { RefreshCwIcon, UserPlusIcon } from 'lucide-svelte';
	import { onMount, type ComponentProps } from 'svelte';
	import { noop } from 'svelte/internal';

	let refreshProgress: Promise<void> = new Promise(noop);
	let lastRefresh: Date = new Date();

	let players: Player[] = [];
	async function fetchPlayers() {
		const rpc = await getAdminServiceClient();
		const resp = await rpc.listPlayers({
			eventKey: $page.params.eventKey
		});

		players = resp.players;
	}

	async function refreshData() {
		refreshProgress = fetchPlayers();
		await refreshProgress;
		lastRefresh = new Date();
	}
	onMount(refreshData);

	function showModal(title: string, playerData: PlayerData, playerId?: string) {
		let operation: FormOperation = 'create';
		if (playerId) {
			operation = 'edit';
		}

		const props: Omit<ComponentProps<PlayerForm>, 'parent'> = {
			operation,
			title,
			playerData,
			onSubmit: async (op: FormOperation, playerData: PlayerData): Promise<DisplayError> => {
				const rpc = await getAdminServiceClient();
				try {
					if (operation === 'create') {
						await rpc.createPlayer({
							eventKey: $page.params.eventKey,
							playerData
						});
					}

					if (operation === 'edit') {
						await rpc.updatePlayer({
							playerId,
							playerData
						});
					}
				} catch (error) {
					console.log('API Error', error);
					return { type: 'Server Error', message: error as string };
				}

				return null;
			}
		};

		modalStore.trigger({
			type: 'component',
			component: {
				ref: PlayerForm,
				props
			},
			response: (submittedForm: boolean) => submittedForm && refreshData()
		});
	}

	function showNewPlayerModal() {
		showModal('Register New Player', new PlayerData());
	}

	function getPlayerCategory(player: Player) {
		const cat = player.events.find((x) => x.eventKey == $page.params.eventKey)?.scoringCategory;
		return cat ? scoringCategoryToDisplayName[cat] : '(None)';
	}
</script>

<SetTitle title="Players" />

<div class="max-w-3xl mx-auto">
	<RefreshHeader
		title="Players"
		refresh={refreshData}
		loadingStatus={refreshProgress}
		{lastRefresh}
	/>

	{#await refreshProgress}
		<div class="card py-12 flex flex-col items-center">
			<p class="mb-4">Loading players...</p>
			<RefreshCwIcon class="animate-spin" />
		</div>
	{:then}
		{#if players.length}
			<div class="table-container">
				<table class="table table-hover">
					<thead>
						<tr>
							<th>Name</th>
							<th>League</th>
							<th class="table-cell-fit">Edit</th>
						</tr>
					</thead>
					<tbody>
						{#each players as player (player.id)}
							<tr>
								<td>{player.data?.name}</td>
								<td>{getPlayerCategory(player)}</td>
								<td class="table-cell-fit">
									<button
										type="button"
										class="btn btn-sm variant-filled"
										on:click={() =>
											showModal('Update Player', player.data || new PlayerData(), player.id)}
									>
										<span>Edit</span>
									</button>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{:else}
			<NoDataCard
				text="No Players to Display"
				ctaText="Register a New Player"
				on:click={showNewPlayerModal}
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
		<button type="button" class="btn btn-lg variant-filled-primary" on:click={showNewPlayerModal}>
			<span class="badge-icon"><UserPlusIcon /></span>
			<span>New Player</span>
		</button>
	</footer>
</div>

<style lang="postcss">
	.table tbody td {
		@apply align-middle;
	}
</style>
