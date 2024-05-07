<script lang="ts">
	import { page } from '$app/stores';
	import ErrorBanner from '$lib/components/ErrorBanner.svelte';
	import NoDataCard from '$lib/components/dashboards/NoDataCard.svelte';
	import RefreshHeader from '$lib/components/dashboards/RefreshHeader.svelte';
	import StageForm from '$lib/components/modals/StageForm.svelte';
	import SetTitle from '$lib/components/util/SetTitle.svelte';
	import { UpdateStageRequest, type Stage } from '$lib/proto/api/v1/admin_pb';
	import type { Venue } from '$lib/proto/api/v1/shared_pb';
	import { getAdminServiceClient } from '$lib/rpc/client';
	import { modalStore } from '@skeletonlabs/skeleton';
	import { RefreshCwIcon } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { noop, type ComponentProps } from 'svelte/internal';

	let refreshProgress: Promise<void> = new Promise(noop);
	let lastRefresh: Date;

	let venues: Venue[] = [];
	async function fetchVenues() {
		const rpc = await getAdminServiceClient();
		const resp = await rpc.listVenues({});

		venues = resp.venues;
	}

	let stages: Stage[] = [];
	async function fetchStages() {
		const rpc = await getAdminServiceClient();
		const resp = await rpc.listEventStages({
			eventKey: $page.params.eventKey
		});

		stages = resp.stages;
	}

	async function refreshData() {
		refreshProgress = Promise.all([fetchVenues(), fetchStages()]).then(noop);
		await refreshProgress;
		lastRefresh = new Date();
	}
	onMount(refreshData);

	async function showEditStageModal(stage: Stage) {
		const props: Omit<ComponentProps<StageForm>, 'parent'> = {
			venues: await venues,
			stage: new UpdateStageRequest({
				stageId: stage.id,
				venueId: stage.venue?.id,
				rank: stage.rank,
				durationMin: stage.durationMin,
				venueDescription: stage.rule?.venueDescription
			}),
			onSubmit: async (data: UpdateStageRequest) => {
				const rpc = await getAdminServiceClient();
				try {
					await rpc.updateStage(data);
				} catch (error) {
					return { type: 'Server Error', message: error as string };
				}

				return null;
			}
		};

		modalStore.trigger({
			type: 'component',
			component: {
				ref: StageForm,
				props
			},
			response: (submittedForm: boolean) => submittedForm && refreshData()
		});
	}
</script>

<SetTitle title="Schedule" />

<div class="max-w-3xl mx-auto mb-4">
	<RefreshHeader
		title="Schedule"
		refresh={refreshData}
		loadingStatus={refreshProgress}
		{lastRefresh}
	/>

	{#await refreshProgress}
		<div class="card py-12 flex flex-col items-center">
			<p class="mb-4">Loading schedule...</p>
			<RefreshCwIcon class="animate-spin" />
		</div>
	{:then}
		{#if stages.length}
			<div class="table-container">
				<table class="table table-hover">
					<thead>
						<tr>
							<th>Venue</th>
							<th>Duration (Min)</th>
							<th class="long-text-col">Rules</th>
							<th class="table-cell-fit">Actions</th>
						</tr>
					</thead>
					<tbody>
						{#each stages as stage (stage.id)}
							<tr>
								<td>{stage.venue?.name}</td>
								<td>{stage.durationMin}</td>
								<td class="long-text-col">{stage.rule?.venueDescription}</td>
								<td class="table-cell-fit action-btns text-right">
									<button
										type="button"
										class="btn btn-sm variant-filled ml-4"
										on:click={() => showEditStageModal(stage)}
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
			<NoDataCard text="No Stops to Display" hideCTA={true} />
		{/if}
	{:catch error}
		<ErrorBanner
			error={{ type: 'Server Error', message: error }}
			dismissLabel="Retry"
			on:dismiss={refreshData}
		/>
	{/await}
</div>

<style lang="postcss">
	.table tbody td {
		@apply align-middle;
	}
	.long-text-col {
		@apply whitespace-normal;
	}
	.action-btns {
		@apply whitespace-nowrap;
	}
</style>
