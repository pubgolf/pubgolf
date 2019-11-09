<script>
  let code = '';

  function submit() {
    console.log('Verifying', code);
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
