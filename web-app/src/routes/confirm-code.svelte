<script>
  let code = '';

  function submit() {
    console.log('Verifying', code);
  }
</script>

<style>
</style>

<form on:submit|preventDefault={submit}>
  <label for="confirm-code">
    Enter the code you received
  </label>
  <input
    id="confirm-code"
    class="input w-full"
    type="tel"
    placeholder="123456"
    bind:value={code}
  >

  <button class="btn btn-primary w-full mt-4">
    Submit
  </button>
</form>


