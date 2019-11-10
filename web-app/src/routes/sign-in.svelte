<script>
  import { goto } from "@sapper/app";

  let phone = '';
  // TODO: format phone as they type
  // TODO: validation

  function submit() {
    console.log(`Sending SMS to ${phone}`);

    goto('confirm-code');
  }
</script>

<svelte:head>
  <title>Sign into Pub Golf</title>
</svelte:head>

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
