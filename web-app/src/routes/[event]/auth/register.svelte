<script>
  import { goto } from '@sapper/app';

  import { LEAGUE } from '../../../api';
  import {
    api,
    event,
  } from '../../../stores';
  import FormError from './_FormError';


  // Local state
  let name = '';
  let phone = '';
  let league = '';

  // reset error to null if the form changes
  $: error = Boolean(name && phone && league) && null;

  function submit () {
    const player = { name, phone, league };

    // console.log('Registering', player);

    error = null;
    $api.registerPlayer(player).then(() => {
      goto(`${$event}/auth/confirm?phone=${phone}`);
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
    bind:value="{phone}"
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
