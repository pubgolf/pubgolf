<script context="module">
  export const preload = (page) => ({ eventKey: page.params.event });
</script>

<script>
  import {
    goto,
    stores,
  } from '@sapper/app';

  import { getAPI } from 'src/api';
  import {
    applyInsertions,
    onlyDigits,
  } from 'src/phone-handler';
  import FormError from './_FormError';

  export let eventKey;

  const { session } = stores();
  const api = getAPI($session);

  // Local state
  let phone = '';

  // reset error to null if the form changes
  $: error = Boolean(phone) && null;

  function handlePhone (inputEvent) {
    const { target } = inputEvent;
    phone = applyInsertions(target.value);

    // puts cursor back to the position where addition or deletion was done
    // otherwise it always jumps back to the end.
    let position = target.selectionEnd;
    const digit = target.value[position - 1];
    target.value = applyInsertions(target.value);
    while (position < target.value.length && target.value.charAt(position - 1) !== digit) {
      position += 1;
    }
    setTimeout(() => {
      target.selectionStart = position;
      target.selectionEnd = position;
    }, 0);
  }

  function submit () {
    // console.log(`Requesting login for ${phone}`);
    const unformattedPhone = onlyDigits(phone);

    error = null;
    api.requestPlayerLogin({
      eventKey,
      phoneNumber: `+1${unformattedPhone}`,
    }).then(() => {
      goto(`${eventKey}/auth/confirm?phone=${unformattedPhone}`);
    }, (apiError) => {
      error = apiError;
    });
  }
</script>

<svelte:head>
  <title>Log In | Pub Golf</title>
</svelte:head>

<FormError {error}/>

<form on:submit|preventDefault="{submit}">
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
    on:input="{handlePhone}"
    required
  >

  <button class="btn btn-primary w-full mt-4">
    Send code
  </button>
</form>
