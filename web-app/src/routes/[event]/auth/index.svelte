<script>
  import { event } from 'src/stores';
</script>

<svelte:head>
  <title>Pub Golf</title>
</svelte:head>

<a
  href="{$event}/auth/login"
  class="btn btn-primary"
>
  Sign In
</a>
<span class="my-2 text-center text-white">or</span>
<a
  href="{$event}/auth/register"
  class="btn btn-secondary"
>
  Register
</a>
