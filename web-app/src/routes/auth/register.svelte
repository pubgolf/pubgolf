<script>
  import { goto } from '@sapper/app';

  import { LEAGUE, DEFAULT_CLIENT } from '../../api';
  import { player } from '../../stores';


  let name = '';
  let phone = '';
  let league = '';

  function submit () {
    $player = { name, phone, league };

    console.log('Registering', $player);

    DEFAULT_CLIENT.registerPlayer($player).then(() => {
      goto('auth/confirm-code');
    });
  }
</script>

<style>
  form {
    display: grid;
    grid-template: /* @formatter:off */
      "label-name  input-name " auto
      "label-phone input-phone" auto
      "league      league-opts" auto
      "submit      submit     " auto /
       1fr         1fr; /* @formatter:on */
    grid-gap: 0.5rem;
  }

  .SUBMIT {
    grid-area: submit;
  }
</style>

<svelte:head>
  <title>Register for Pub Golf</title>
</svelte:head>

<form on:submit|preventDefault="{submit}">
  <label for="register-name" class="font-bold">
    Name:
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

  <label for="register-phone" class="font-bold">
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

  <span class="font-bold p-1">League:</span>
  <div class="flex">
    <!--  TODO: give these an empty state  -->
    <label class="flex-grow input text-center text-blue-400 mr-2">
      <input
        type="radio"
        name="league"
        value="{LEAGUE.PGA}"
        bind:group="{league}"
      >
      PGA
    </label>
    <label class="flex-grow input text-center text-blue-400">
      <input
        type="radio"
        name="league"
        value="{LEAGUE.LPGA}"
        bind:group="{league}"
      >
      LPGA
    </label>
  </div>

  <button class="SUBMIT btn btn-primary mt-2">
    Register
  </button>
</form>
