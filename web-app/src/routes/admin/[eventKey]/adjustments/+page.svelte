<script lang="ts">
	import { page } from '$app/stores';
	import ErrorBanner from '$lib/components/ErrorBanner.svelte';
	import NoDataCard from '$lib/components/dashboards/NoDataCard.svelte';
	import RefreshHeader from '$lib/components/dashboards/RefreshHeader.svelte';
	import AdjustmentTemplateForm from '$lib/components/modals/AdjustmentTemplateForm.svelte';
	import SetTitle from '$lib/components/util/SetTitle.svelte';
	import type {
		AdjustmentTemplate,
		AdjustmentTemplateData,
		Stage
	} from '$lib/proto/api/v1/admin_pb';
	import { getAdminServiceClient } from '$lib/rpc/client';
	import { modalStore, toastStore } from '@skeletonlabs/skeleton';
	import { RefreshCwIcon, TrophyIcon } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { noop, type ComponentProps } from 'svelte/internal';

	let refreshProgress: Promise<void> = new Promise(noop);
	let lastRefresh: Date;

	let stages: Stage[] = [];
	async function fetchStages() {
		const rpc = await getAdminServiceClient();
		const resp = await rpc.listEventStages({
			eventKey: $page.params.eventKey
		});

		stages = resp.stages;
	}

	let templates: AdjustmentTemplate[] = [];
	async function fetchTemplates() {
		const rpc = await getAdminServiceClient();
		const resp = await rpc.listAdjustmentTemplates({
			eventKey: $page.params.eventKey
		});

		templates = resp.templates;
	}

	async function refreshData() {
		refreshProgress = Promise.all([fetchStages(), fetchTemplates()]).then(noop);
		await refreshProgress;
		lastRefresh = new Date();
	}
	onMount(refreshData);

	function getVenueName(stageId?: string) {
		if (!stageId) {
			return 'All Venues';
		}
		const i = stages.findIndex((x) => x.id === stageId);
		return `${i + 1}: ${stages[i].venue?.name}`;
	}

	async function setTemplateVisibility(template: AdjustmentTemplate, visible: boolean) {
		let action = 'Hide';
		if (visible) {
			action = 'Un-Hide';
		}

		modalStore.trigger({
			type: 'confirm',
			title: `Confirm ${action}`,
			body: `Are you sure you wish to ${action.toLocaleLowerCase()} this template ("${
				template.data?.adjustment?.label
			}" on "${getVenueName(template.data?.stageId)}") from players?`,
			response: async (r: boolean) => {
				if (!r || !template.data) {
					return;
				}
				template.data.isVisible = visible;
				const rpc = await getAdminServiceClient();
				try {
					await rpc.updateAdjustmentTemplate({ template });
				} catch (error) {
					toastStore.trigger({
						message: `API Error: ${error}`,
						background: 'variant-filled-error'
					});
				}

				refreshData();
			}
		});
	}

	async function showNewTemplateModal() {
		const props: Omit<ComponentProps<AdjustmentTemplateForm>, 'parent'> = {
			operation: 'create',
			eventKey: $page.params.eventKey,
			stages: await stages,
			template: null,
			onSubmit: async (data: AdjustmentTemplateData, id?: string) => {
				const rpc = await getAdminServiceClient();
				try {
					await rpc.createAdjustmentTemplate({
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
				ref: AdjustmentTemplateForm,
				props
			},
			response: (submittedForm: boolean) => submittedForm && refreshData()
		});
	}

	async function showEditTemplateModal(template: AdjustmentTemplate) {
		const props: Omit<ComponentProps<AdjustmentTemplateForm>, 'parent'> = {
			operation: 'edit',
			eventKey: $page.params.eventKey,
			stages: await stages,
			template,
			onSubmit: async (data: AdjustmentTemplateData, id?: string) => {
				if (!id) {
					return { type: 'App Error', message: 'Did not receive ID for template' };
				}

				const rpc = await getAdminServiceClient();
				try {
					await rpc.updateAdjustmentTemplate({ template: { id, data } });
				} catch (error) {
					return { type: 'Server Error', message: error as string };
				}

				return null;
			}
		};

		modalStore.trigger({
			type: 'component',
			component: {
				ref: AdjustmentTemplateForm,
				props
			},
			response: (submittedForm: boolean) => submittedForm && refreshData()
		});
	}
</script>

<SetTitle title="Adjustments" />

<div class="max-w-3xl mx-auto mb-4">
	<RefreshHeader
		title="Adjustments"
		refresh={refreshData}
		loadingStatus={refreshProgress}
		{lastRefresh}
	/>

	{#await refreshProgress}
		<div class="card py-12 flex flex-col items-center">
			<p class="mb-4">Loading adjustment templates...</p>
			<RefreshCwIcon class="animate-spin" />
		</div>
	{:then}
		{#if templates.length}
			<div class="table-container">
				<table class="table table-hover">
					<thead>
						<tr>
							<th>Venue</th>
							<th>Label</th>
							<th>Value</th>
							<th>Rank</th>
							<th class="table-cell-fit">Actions</th>
						</tr>
					</thead>
					<tbody>
						{#each templates as template (template.id)}
							<tr>
								<td>{getVenueName(template.data?.stageId)}</td>
								<td>{template.data?.adjustment?.label}</td>
								<td>{template.data?.adjustment?.value}</td>
								<td>{template.data?.rank}</td>
								<td class="table-cell-fit action-btns text-right">
									<button
										type="button"
										class="btn btn-sm variant-filled ml-4"
										on:click={() => showEditTemplateModal(template)}
									>
										<span>Edit</span>
									</button>
									{#if template.data?.isVisible}
										<button
											type="button"
											class="btn btn-sm variant-filled-error ml-4"
											on:click={() => setTemplateVisibility(template, false)}
										>
											<span>Hide</span>
										</button>
									{:else}
										<button
											type="button"
											class="btn btn-sm variant-outline-error ml-4"
											on:click={() => setTemplateVisibility(template, true)}
										>
											<span>Show</span>
										</button>
									{/if}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{:else}
			<NoDataCard
				text="No Adjustments to Display"
				ctaText="Create an Adjustment"
				on:click={showNewTemplateModal}
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
		<button type="button" class="btn btn-lg variant-filled-primary" on:click={showNewTemplateModal}>
			<span class="badge-icon"><TrophyIcon /></span>
			<span>Create Adjustment</span>
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
