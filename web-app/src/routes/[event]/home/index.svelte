<script>
  import { goto } from '@sapper/app';

  import Countdown from '../../../components/Countdown';
  import {
    api,
    event,
    stops,
    nextStop,
    pastStops,
  } from '../../../stores';

  // If not authenticated, this page isn't accessible
  const redirectIfUnauthorized = () => {
    if (!$api.isLoggedIn()) {
      goto(`${$event}/auth`);
    }
  };

  let fetching = true;

  $: if ($event) {
    if (!$api.isLoggedIn()) {
      goto(`${$event}/auth`);
    } else {
      $api.getSchedule().then((schedule) => {
        $stops = schedule.venuelist.venuesList.map(
          ({ stopid, venue: { starttime, ...venue } }, index) => ({
              stopid,
              ...venue,
              start: new Date(starttime),
              index,
            }
          ));
        fetching = false;
      }, () => {
        if (!$api.isLoggedIn()) {
          goto(`${$event}/auth`);
        }
      });
    }
  }

  // $: console.log($nextStop);
</script>

<style>
  .HOME {
    height: 90%;
  }
</style>

{#if fetching}
  <div class="flex h-full items-center text-center text-6xl">
    Fetching Schedule
  </div>
{:else}
  <div class="HOME text-center text-4xl pt-32">
    {#if $nextStop}
      <p>
        Time Remaining
      </p>
      <p class="text-red-600">
        <Countdown to="{$nextStop.start}"/>
      </p>
      <p class="text-2xl mt-32">
        Next Bar: {$nextStop.name}<br>
        <a
          href="https://www.google.com/maps/place/{$nextStop.address}"
          target="_blank"
          class="text-blue-400 underline text-xl"
        >
          {$nextStop.address}
        </a>
      </p>
      <a
        href="{$event}/home/add-score"
        class="block btn btn-primary w-2/3 my-16 mx-auto"
      >
        Add your Score
      </a>
    {:else if $stops.length}
      <p class="text-6xl">
        Thanks for Playing!
      </p>
    {/if}
  </div>

  {#if $pastStops.length}
    <ol class="border-t-4 border-gray-500 px-2">
      {#each $pastStops as stop, i (stop.stopid)}
        <li class="{ i ? 'border-t-2 ' : ''}border-gray-500 px-1">
          {#if i === 0 && $nextStop !== null}
            <span class="text-green-500 italic">Current:</span>
          {/if}
          {stop.name}<br>
          <a
            href="https://www.google.com/maps/place/{stop.address}"
            target="_blank"
            class="text-blue-400 underline"
          >
            {stop.address}
          </a>
        </li>
      {/each}
    </ol>
  {/if}
{/if}

