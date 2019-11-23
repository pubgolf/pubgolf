<script context="module">
  export function preload (page) {
    this.redirect(301, `${page.path}/auth`);
  }

  // Looking for where the API gets initialized?
  // Go to web-app/src/routes/[event]/_layout
</script>
