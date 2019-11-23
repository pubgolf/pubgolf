<script>
  export let basePath;
  export let links;
  export let segment;

  $: console.log(segment);
</script>

<style>
  .LIST {
    display: flex;
    margin: 0;
    padding: 0;
    -webkit-box-shadow: 0px -5px 10px 0px rgba(136, 136, 136, 0.3);
    -moz-box-shadow: 0px -5px 10px 0px rgba(136, 136, 136, 0.3);
    box-shadow: 0px -5px 10px 0px rgba(136, 136, 136, 0.3);
  }

  /* clearfix */
  .LIST:after {
    content: '';
    display: block;
    clear: both;
  }

  .LIST > * {
    flex-grow: 1;
  }

  .LINK {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }

  .ICON {
    width: 1.8rem;
    margin: 0.25rem 0 0;
  }
</style>

{#if links && links.length}
  <nav class="flex-grow-0 flex-shrink-0 bg-gray-200 font-light">
    <ul class="LIST">
      {#each links as link (link.segment)}
        <li>
          <a
            class="LINK{segment === link.segment ? ' bg-orange text-white' : ''}"
            href="{basePath}/{link.segment}"
          >
            <div class="ICON">{@html link.icon}</div>
            {link.text}
          </a>
        </li>
      {/each}
    </ul>
  </nav>
{/if}
