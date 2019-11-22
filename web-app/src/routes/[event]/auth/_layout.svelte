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
    grid-gap: 0.5em;
    height: 100%;
    background-color: #50AF4F;
    color: white;
    padding: 1em;
    margin: 0 auto;
  }

  .IMAGE {
    width: 100%;
    max-width: 250px;
    margin: 0 auto;
  }

  .HEADING {
    font-size: 3em;
    text-align: center;
    font-weight: 700;
    margin: 0;
  }
</style>

<main class="MAIN">
  <img class="IMAGE" alt="Cheers" src="social-beer--green.png">

  <h1 class="HEADING">Pub Golf NYC 2k19</h1>

  <div class="flex flex-col w-full text-xl mt-8 mx-auto">
    <slot></slot>
  </div>
</main>
