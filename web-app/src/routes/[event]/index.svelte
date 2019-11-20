<script context="module">
  export function preload (page) {
    this.redirect(301, `${page.path}/auth`);
  }
</script>
