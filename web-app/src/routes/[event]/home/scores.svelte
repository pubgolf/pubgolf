<script context="module">
  export async function preload () {
    console.info('scores preload');
  }
</script>

<script>
  import { stores } from '@sapper/app';
  import { onMount } from 'svelte';

  import {
    api,
    event,
  } from 'src/stores';

  let scores;
  const { page, session } = stores();

  onMount(async () => {
    try {
      scores = await $api.getScores();
    } catch (e) {
      console.error(e);
    }
  });

  $: scores && console.log(scores.scorelistsList);
</script>

<style>
</style>
