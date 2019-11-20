<script>
  export let error;
</script>

{#if error}
  <p class="text-red-600 p-1 mb-2">
    <span class="font-extrabold uppercase">Error:</span>
    {error.message}
  </p>
{/if}
