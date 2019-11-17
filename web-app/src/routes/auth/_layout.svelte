<script>
  // TODO: transition between pages
  export let segment;
</script>

<style>
  .MAIN {
    display: grid;
    grid-template: /* @formatter:off */
      "image  " min-content
      "heading" min-content
      "content" 1fr /
       1fr; /* @formatter:on */
    height: 100%;
    background-color: white;
    padding: 2em;
    margin: 0 auto;
  }

  .HEADING {
    font-size: 2.8em;
    text-transform: uppercase;
    font-weight: 700;
    margin: 0 0 0.5em 0;
  }

  .IMAGE {
    width: 100%;
    max-width: 400px;
    margin: 0 0 1em 0;
  }
</style>

<main class="MAIN">
  <img class="IMAGE" alt="Borat" src="great-success.png">

  <h1 class="HEADING">Pub Golf</h1>

  <div class="flex flex-col w-full mx-auto">
    <slot></slot>
  </div>
</main>
