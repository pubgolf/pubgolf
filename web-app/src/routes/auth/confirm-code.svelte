<script>
  import { goto } from '@sapper/app';

  import { DEFAULT_CLIENT } from '../../api';
  import { player } from '../../stores';


  // TODO: There should be a way to pass this through the URL instead
  if (!$player || !$player.phone) {
    goto('auth');
  }

  let code = '';

  function submit () {
    console.log('Verifying', code);

    DEFAULT_CLIENT.playerLogin(
      $player.phone,
      Number(code),
    ).then((response) => {
      debugger;
      goto('app');
    });
  }
</script>

<style>
</style>

<form on:submit|preventDefault="{submit}" class="w-2/3 mx-auto">
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
