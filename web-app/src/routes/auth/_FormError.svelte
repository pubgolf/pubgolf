<script>
  export let error;

  $: if (error && process.env.NODE_ENV === 'development') {
    console.error(error);
  }
</script>

{#if error}
  <p class="text-red-600 p-1 mb-2">
    <span class="font-extrabold uppercase">Error:</span>
    {error.message}
  </p>
{/if}
