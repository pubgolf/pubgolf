<script context="module">
  export function preload (page) {
    this.redirect(302, `${page.params.event}/auth`);
  }

  // Looking for where the API gets initialized?
  // Go to web-app/src/routes/[event]/_layout
</script>
