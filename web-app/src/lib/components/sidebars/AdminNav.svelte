<script lang="ts">
	import { page } from '$app/stores';
	import { drawerStore } from '@skeletonlabs/skeleton';

	const dashboards = [
		{
			name: 'event-info',
			icon: '📝'
		},
		// {
		// 	name: 'schedule',
		// 	icon: '⏱️'
		// },
		{
			name: 'players',
			icon: '🏌️'
		},
		{
			name: 'scores',
			icon: '🏆'
		},
		{
			name: 'alerts',
			icon: '🚨'
		}
	] as const;

	function slugToTitle(s: string) {
		return s
			.split('-')
			.map((x) => x.charAt(0).toUpperCase() + x.slice(1))
			.join(' ');
	}
</script>

<div class="py-8 px-4">
	<span class="h2 ml-4 mb-4">Dashboards</span>
	<nav class="list-nav">
		<ul>
			{#each dashboards as dashboard}
				<li>
					<a
						class:bg-primary-500={$page.route.id?.endsWith(dashboard.name)}
						on:click={() => drawerStore.close()}
						href="../{dashboard.name}/"
					>
						<span class="badge-icon">{dashboard.icon}</span>
						<span class="flex-auto">{slugToTitle(dashboard.name)}</span>
					</a>
				</li>
			{/each}
		</ul>
	</nav>
</div>
