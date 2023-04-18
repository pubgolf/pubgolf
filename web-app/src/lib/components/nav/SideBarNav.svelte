<script lang="ts">
	import { page } from '$app/stores';
	import { drawerStore } from '@skeletonlabs/skeleton';

	export let title: string;
	export let items: {
		slug: string;
		icon: string;
		title?: string;
	}[];
	export let floatOnDesktop = false;

	function slugToTitle(s: string) {
		return s
			.split('-')
			.map((x) => x.charAt(0).toUpperCase() + x.slice(1))
			.join(' ');
	}
</script>

<!-- Touch file -->

<div class="py-8 px-4" class:floating-sidebar={floatOnDesktop}>
	<span class="title h2 ml-4 mb-4">{title}</span>
	<nav class="list-nav">
		<ul>
			{#each items as item}
				<li>
					<a
						class:bg-primary-500={$page.route.id?.endsWith(item.slug)}
						on:click={() => drawerStore.close()}
						href="../{item.slug}/"
					>
						<span class="badge-icon shadow-none">{item.icon}</span>
						<span class="flex-auto">{item.title || slugToTitle(item.slug)}</span>
					</a>
				</li>
			{/each}
		</ul>
	</nav>
</div>

<style lang="postcss">
	.floating-sidebar .title {
		@apply text-2xl md:text-xl md:mb-2 !important;
	}
</style>
