<script lang="ts">
	export let title: string;
	let pageTitle = title ? `${title} | PubGolf.co` : 'PubGolf.co';
</script>

<svelte:head>
	<title>{pageTitle}</title>
</svelte:head>
