<script lang="ts">
	import { modalStore } from '@skeletonlabs/skeleton';
	import { RefreshCcwIcon, UserPlusIcon } from 'lucide-svelte';
	import type { ComponentProps } from 'svelte';
	import { AdminClient } from '$lib/rpc/client';
	import { page } from '$app/stores';
	import SetTitle from '$lib/components/util/SetTitle.svelte';
	import PlayerForm, {
		type FormError,
		type FormOperation
	} from '$lib/components/modals/PlayerForm.svelte';
	import { PlayerData, type Player } from '$lib/proto/api/v1/shared_pb';
	import { scoringCategoryToDisplayName } from '$lib/models/scoring-category';

	let players: Promise<Player[]> = fetchPlayers();
	let dataUpdatedAt: Date = new Date();

	async function fetchPlayers() {
		const resp = await AdminClient.listPlayers({
			eventKey: $page.params.eventKey
		});

		dataUpdatedAt = new Date();

		return resp.players;
	}

	function refreshData() {
		players = fetchPlayers();
	}

	function showModal(title: string, playerData: PlayerData, playerId?: string) {
		let operation: FormOperation = 'create';
		if (playerId) {
			operation = 'edit';
		}

		const props: Omit<ComponentProps<PlayerForm>, 'parent'> = {
			operation,
			playerData,
			onSubmit: async (op: FormOperation, playerData: PlayerData): Promise<FormError> => {
				try {
					if (operation === 'create') {
						await AdminClient.createPlayer({
							eventKey: $page.params.eventKey,
							playerData
						});
					}

					if (operation === 'edit') {
						await AdminClient.updatePlayer({
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
			title,
			component: {
				ref: PlayerForm,
				props
			},
			response: (submittedForm: boolean) => {
				if (submittedForm) {
					refreshData();
				}
			}
		});
	}

	function showNewPlayerModal() {
		showModal('Register New Player', new PlayerData());
	}
</script>

<SetTitle title="Players" />

<div class="max-w-3xl mx-auto">
	<header class="flex justify-between items-start mb-4 md:mt-4">
		<h1>Players</h1>
		<div class="text-right">
			<button type="button" class="btn variant-filled mb-0.5" on:click={refreshData}>
				<span class="badge-icon"><RefreshCcwIcon /></span>
				<span>Refresh</span>
			</button><br />
			<span class="text-xs"
				>Last Fetched: {new Intl.DateTimeFormat(undefined, { timeStyle: 'medium' }).format(
					dataUpdatedAt
				)}</span
			>
		</div>
	</header>

	{#await players}
		<p>Fetching players...</p>
	{:then players}
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
					{#if players.length}
						{#each players as player (player.id)}
							<tr>
								<td>{player.data?.name}</td>
								<td
									>{player.data?.scoringCategory
										? scoringCategoryToDisplayName[player.data?.scoringCategory]
										: '(None)'}</td
								>
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
					{:else}
						<tr class="text-center">
							<td colspan="3">
								<div class="py-12">
									<p class="mb-4">No Players</p>
									<button
										type="button"
										class="btn btn-lg variant-filled-secondary"
										on:click={showNewPlayerModal}
									>
										<span>Register a New Player</span>
									</button>
								</div>
							</td>
						</tr>
					{/if}
				</tbody>
			</table>
		</div>
	{:catch error}
		<p>Error fetching players: {error}</p>
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
