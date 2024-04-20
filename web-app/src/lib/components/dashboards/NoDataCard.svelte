<script lang="ts">
	export let text = 'No Data to Display';
	export let ctaText = 'Create';
	export let hideCTA = false;
</script>

<div class="card py-12 flex flex-col items-center">
	<p class="mb-4">{text}</p>
	{#if !hideCTA}
		<button type="button" class="btn btn-lg variant-filled-secondary" on:click>
			<span>{ctaText}</span>
		</button>
	{/if}
</div>
