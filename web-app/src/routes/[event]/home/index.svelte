<script>
  import { onMount } from 'svelte';
  import { goto } from '@sapper/app';

  import Countdown from '../../../components/Countdown';
  import {
    api,
    event,
    time,
  } from '../../../stores';
  import { dateDelta } from '../../../utils';


  let stops = [];

  $: if ($event) {

    // If not authenticated, this page isn't accessible
    if (!$api.isLoggedIn()) {
      goto(`${$event}/auth`);
    }

    $api.getSchedule().then((schedule) => {
      stops = schedule.venuelist.venuesList.map(
        ({ stopid, venue: { starttime, ...venue } }, index) => ({
            stopid,
            ...venue,
            start: new Date(starttime),
            index,
          }
        ));
    });
  }

  $: nextBar = stops.length
    ? stops.find(stop => stop.start > $time)
    : null;
  $: pastBars = nextBar ? stops.slice(0, nextBar.index).reverse() : [];

  $: console.log(nextBar, pastBars.length);
</script>

<style>
  .HOME {
    height: 90%;
  }
</style>

<div class="HOME text-center text-4xl pt-32">
  {#if nextBar}
    <p>
      Time Remaining
    </p>
    <p class="text-red-600">
      <Countdown to="{nextBar.start}"/>
    </p>
    <p class="text-2xl mt-32">
      Next Bar: {nextBar.name}<br>
      <a
        href="https://www.google.com/maps/place/{nextBar.address}"
        target="_blank"
        class="text-blue-400 underline text-xl"
      >
        {nextBar.address}
      </a>
    </p>
  {/if}
</div>

<ol class="border-t-4 border-gray-500 px-2">
  {#each pastBars as bar, i}
    <li class="{ i ? 'border-t-2 ' : ''}border-gray-500 px-1">
      {#if i === 0}<span class="text-green-500 italic">Current:</span>{/if}
      {bar.name}<br>
      <a
        href="https://www.google.com/maps/place/{bar.address}"
        target="_blank"
        class="text-blue-400 underline"
      >
        {bar.address}
      </a>
    </li>
  {/each}
</ol>
