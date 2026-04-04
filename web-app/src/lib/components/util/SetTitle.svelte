<script lang="ts">
	let { title }: { title: string } = $props();
	let pageTitle = $derived(title ? `${title} | PubGolf.co` : 'PubGolf.co');
</script>

<svelte:head>
	<title>{pageTitle}</title>
</svelte:head>
