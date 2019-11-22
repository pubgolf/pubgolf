<script>
  import { goto } from '@sapper/app';

  import {
    api,
    event,
  } from '../../../stores';
  import FormError from './_FormError';

  // TODO: format phone as they type
  // TODO: validation

  // Local state
  let phone = '';

  // reset error to null if the form changes
  $: error = Boolean(phone) && null;

  function submit () {
    console.log(`Requesting login for ${phone}`);

    error = null;
    $api.requestPlayerLogin(phone).then(() => {
      // TODO: figure out how to get relative urls
      goto(`${$event}/auth/confirm?phone=${phone}`);
    }, (apiError) => {
      error = apiError;
    });
  }
</script>

<svelte:head>
  <title>Log In | Pub Golf</title>
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
