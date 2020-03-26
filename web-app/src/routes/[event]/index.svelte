<script context="module">
  import { authHelper } from 'src/auth-helper';

  export function preload (page, session) {
    const eventKey = page.params.event;
    if (authHelper.isAuthorized(session, eventKey)) {
      this.redirect(302, `${eventKey}/home`);
    } else {
      this.redirect(302, `${eventKey}/auth`);
    }
  }
</script>
