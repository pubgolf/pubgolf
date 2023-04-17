<script>
	import { page } from '$app/stores';
	import { drawerStore } from '@skeletonlabs/skeleton';
</script>

<div class="py-8 px-4">
	<span class="h2 title">About PubGolf</span>
	<nav class="list-nav">
		<ul>
			<li>
				<a
					href="/about/privacy/"
					class:bg-primary-500={$page.route.id?.endsWith('privacy')}
					on:click={() => drawerStore.close()}
				>
					<span class="badge-icon">🕵️</span>
					<span class="flex-auto">Privacy Policy</span>
				</a>
				<a
					href="/about/contact/"
					class:bg-primary-500={$page.route.id?.endsWith('contact')}
					on:click={() => drawerStore.close()}
				>
					<span class="badge-icon">📞</span>
					<span class="flex-auto">Contact Us</span>
				</a>
			</li>
		</ul>
	</nav>
</div>

<style lang="postcss">
	.title {
		@apply ml-4 text-2xl md:text-xl mb-4 md:mb-2 !important;
	}
</style>
