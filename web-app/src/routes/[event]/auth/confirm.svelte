<script context="module">
  export async function preload (page) {
    if (!page.query.phone) {
      this.redirect(302, `${page.params.event}/auth`);
    }

    return {
      phone: page.query.phone,
    };
  }
</script>


<script>
  import { goto } from '@sapper/app';

  import {
    api,
    event,
  } from '../../../stores';
  import FormError from './_FormError';


  // props
  export let phone;

  // Local state
  let code = '';

  // reset error to null if the form changes
  $: error = Boolean(code) && null;

  function submit () {
    // console.log('Verifying', code);

    error = null;
    $api.playerLogin(phone, Number(code))
      .then(() => {
        goto(`${$event}/home`);
      }, (apiError) => {
        error = apiError;
      });
  }
</script>

<style>
</style>

<FormError {error}/>

<form on:submit|preventDefault="{submit}">
  <label for="confirm-code">
    Enter the code you received
  </label>
  <input
    id="confirm-code"
    class="input w-full"
    type="tel"
    autocomplete="none"
    placeholder="123456"
    bind:value="{code}"
    required
  >

  <button class="btn btn-primary w-full mt-4">
    Submit
  </button>
</form>
