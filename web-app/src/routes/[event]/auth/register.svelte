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
  let name = '';
  let phone = '';
  const league = '';

  // reset error to null if the form changes
  $: error = Boolean(name && phone && league) && null;

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
    // console.log('Registering', player);
    const unformattedPhone = onlyDigits(phone);

    error = null;
    api.registerPlayer({
      name,
      phoneNumber: `+1${unformattedPhone}`,
      league,
      eventKey,
    }).then(() => {
      goto(`${eventKey}/auth/confirm?phone=${unformattedPhone}`);
    }, (apiError) => {
      error = apiError;
    });
  }
</script>

<style>
  form {
    display: grid;
    grid-template: /* @formatter:off */
      "label-name  input-name " auto
      "label-phone input-phone" auto
      /*"league      league-opts" auto*/
      "submit      submit     " auto /
       1fr         1fr; /* @formatter:on */
    grid-gap: 0.5rem;
    align-items: baseline;
  }

  .SUBMIT {
    grid-area: submit;
  }
</style>

<svelte:head>
  <title>Register | Pub Golf</title>
</svelte:head>

<FormError {error}/>

<form on:submit|preventDefault="{submit}">
  <label for="register-name" class="text-2xl">
    Full Name:
  </label>
  <input
    id="register-name"
    class="input w-full mb-2"
    type="text"
    name="name"
    autocomplete="name"
    placeholder="Full Name"
    bind:value="{name}"
    required
  >

  <label for="register-phone" class="text-2xl">
    Mobile Phone:
  </label>
  <input
    id="register-phone"
    class="input w-full mb-2"
    type="tel"
    name="phone"
    autocomplete="tel"
    placeholder="(123) 555-1234"
    maxlength="14"
    bind:value="{phone}"
    on:input="{handlePhone}"
    required
  >

  <!--<span class="text-2xl">League:</span>
  <div class="flex">
    &lt;!&ndash;  TODO: give these an empty state  &ndash;&gt;
    <label class="flex-grow input text-center text-orange mr-2">
      <input
        type="radio"
        name="league"
        value="{LEAGUE.PGA}"
        bind:group="{league}"
        required
      >
      PGA
    </label>
    <label class="flex-grow input text-center text-orange">
      <input
        type="radio"
        name="league"
        value="{LEAGUE.LPGA}"
        bind:group="{league}"
        required
      >
      LPGA
    </label>
  </div>-->

  <button class="SUBMIT btn btn-primary mt-2">
    Register
  </button>
</form>
