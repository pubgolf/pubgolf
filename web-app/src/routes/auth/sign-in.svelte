<script>
  import { goto } from '@sapper/app';

  import { DEFAULT_CLIENT } from '../../api';
  import { player } from '../../stores';
  import FormError from './_FormError';


  let phone = '';
  // TODO: format phone as they type
  // TODO: validation

  // reset error to null if the form changes
  $: error = Boolean(phone) && null;

  function submit () {
    $player = { phone };

    console.log(`Requesting login for ${phone}`);

    error = null;
    DEFAULT_CLIENT.requestPlayerLogin(phone).then(() => {
      goto('auth/confirm-code');
    }, (apiError) => {
      error = apiError;
    });
  }
</script>

<svelte:head>
  <title>Log in | Pub Golf</title>
</svelte:head>

<FormError {error}/>

<form on:submit|preventDefault="{submit}" class="w-2/3 mx-auto">
  <label for="signin-phone">
    Enter your mobile number
  </label>
  <input
    id="signin-phone"
    class="input w-full"
    type="tel"
    name="phone"
    autocomplete="tel"
    placeholder="(123) 555-1234"
    bind:value="{phone}"
    required
  >

  <button class="btn btn-primary w-full mt-4">
    Send code
  </button>
</form>
