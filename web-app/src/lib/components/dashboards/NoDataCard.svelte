<script lang="ts">
	export let text = 'No Data to Display';
	export let ctaText = 'Create';
</script>

<div class="card py-12 flex flex-col items-center">
	<p class="mb-4">{text}</p>
	<button type="button" class="btn btn-lg variant-filled-secondary" on:click>
		<span>{ctaText}</span>
	</button>
</div>
