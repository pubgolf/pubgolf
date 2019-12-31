<script context="module">
  export const preload = (page) => ({ eventKey: page.params.event });
</script>

<script>
  export let eventKey;
</script>

<svelte:head>
  <title>Pub Golf</title>
</svelte:head>

<a
  href="{eventKey}/auth/login"
  class="btn btn-primary"
>
  Sign In
</a>
<span class="my-2 text-center text-white">or</span>
<a
  href="{eventKey}/auth/register"
  class="btn btn-secondary"
>
  Register
</a>
