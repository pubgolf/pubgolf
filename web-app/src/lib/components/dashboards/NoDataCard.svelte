<script lang="ts">
	interface Props {
		text?: string;
		ctaText?: string;
		hideCTA?: boolean;
		onclick?: (e: MouseEvent) => void;
	}

	let {
		text = 'No Data to Display',
		ctaText = 'Create',
		hideCTA = false,
		onclick
	}: Props = $props();
</script>

<div class="card py-12 flex flex-col items-center">
	<p class="mb-4">{text}</p>
	{#if !hideCTA}
		<button type="button" class="btn btn-lg variant-filled-secondary" {onclick}>
			<span>{ctaText}</span>
		</button>
	{/if}
</div>
