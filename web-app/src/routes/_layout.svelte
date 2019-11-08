<script>
  // TODO: transition between pages
</script>

<style>
  main {
    display: grid;
    grid-template:
      "image  " min-content
      "heading" 1fr
      "content" 50% /
       1fr;
    height: 100%;
    background-color: white;
    padding: 2em;
    margin: 0 auto;
  }
</style>

<main>
  <img alt="Borat" src="great-success.png">

  <h1 class="text-center">Pub Golf</h1>

  <div class="flex flex-col w-2/3 mx-auto">
    <slot></slot>
  </div>
</main>