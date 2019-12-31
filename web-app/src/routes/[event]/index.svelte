<script context="module">
  export function preload (page, session) {
    if (session.user.token) {
      this.redirect(302, `${page.params.event}/home`);
    } else {
      this.redirect(302, `${page.params.event}/auth`);
    }
  }
</script>
