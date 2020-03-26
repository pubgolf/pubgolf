<script context="module">
  import { onlyDigits } from 'src/phone-handler';


  export async function preload (page) {
    if (!/^\d{10}$/.test(page.query.phone)) {
      this.redirect(302, `${page.params.event}/auth`);
    }

    return {
      eventKey: page.params.event,
      phone: onlyDigits(page.query.phone),
    };
  }
</script>


<script>
  import {
    goto,
    stores,
  } from '@sapper/app';
  import { onMount } from 'svelte';

  import { getAPI } from 'src/api';
  import { authHelper } from 'src/auth-helper';
  import { isDev } from 'src/utils';
  import FormError from './_FormError';


  // props
  export let eventKey;
  export let phone;

  const { session } = stores();
  const api = getAPI($session);

  // Local state
  let code = '';

  // reset error to null if the form changes
  $: error = Boolean(code) && null;

  async function resendCode () {
    return await api.requestPlayerLogin({
      eventKey,
      phoneNumber: `+1${phone}`,
    });
  }

  if (isDev($session)) {
    onMount(resendCode);
  }

  function submit () {
    // console.log('Verifying', code);

    error = null;

    api.playerLogin({
      eventKey,
      // TODO: this phone number normalization should be centralized
      phoneNumber: `+1${phone}`,
      authCode: Number(code),
    }).then(async (user) => {
      await authHelper.preserveSession({
        ...user,
        eventKey,
      });
      goto(`${eventKey}/home`);
    }, (apiError) => {
      error = apiError;
    });
  }
</script>

<style>
</style>

<svelte:head>
  <title>Confirm your Phone Number | Pub Golf</title>
</svelte:head>

<FormError {error}/>

<form on:submit|preventDefault="{submit}">
  <label for="confirm-code">
    Enter the code you received
  </label>
  <input
    id="confirm-code"
    class="input w-full"
    type="tel"
    autocomplete="off"
    placeholder="123456"
    bind:value="{code}"
    required
  >

  <button class="btn btn-primary w-full mt-4">
    Submit
  </button>
</form>
