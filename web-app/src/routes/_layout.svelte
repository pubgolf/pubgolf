<script>
  // TODO: transition between pages
</script>

<style>
  main {
    display: grid;
    grid-template:
      'image  ' min-content
      'heading' 1fr
      'content' 50% /
      1fr;
    height: 100%;
    background-color: white;
    padding: 2em;
    margin: 0 auto;
  }

  h1 {
    font-size: 2.8em;
    text-transform: uppercase;
    font-weight: 700;
    margin: 0 0 0.5em 0;
  }

  img {
    width: 100%;
    max-width: 400px;
    margin: 0 0 1em 0;
  }
</style>

<main>
  <img alt="Borat" src="great-success.png" />

  <h1 class="text-center">Pub Golf</h1>

  <div class="flex flex-col w-full mx-auto">
    <slot />
  </div>
</main>
