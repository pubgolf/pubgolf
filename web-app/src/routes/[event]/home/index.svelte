<script context="module">
  import { authHelper } from 'src/auth-helper';

  export async function preload (page, session) {
    if (!authHelper.isAuthorized(session)) {
      try {
        session.user = await authHelper.restoreSession(this.fetch);
      } catch (e) {
        console.debug(e);
        // If not authenticated, this page isn't accessible
        this.redirect(302, `${page.params.event}/auth`);
      }
    }

    return {
      eventKey: page.params.event,
    };
  }
</script>

<script>
  import { onMount } from 'svelte';
  import {
    goto,
    stores,
  } from '@sapper/app';

  import { getAPI } from 'src/api';
  import Countdown from 'src/components/Countdown';
  import {
    stops,
    nextStop,
    pastStops,
  } from 'src/stores';


  // props
  export let eventKey;

  let fetching = true;
  const { page, session } = stores();

  /**
   * Flatten stop and venue into a single object.
   * @param stopid
   * @param venue
   * @param index
   * @returns {{start: Date, stopid: *, index: *}}
   */
  const flattenStop = ({ stopid, venue }, index) => ({
    stopid,
    ...venue,
    start: new Date(venue.starttime),
    index,
  });

  onMount(async () => {
    try {
      const api = getAPI($session);
      const schedule = await api.getSchedule({ eventKey });

      $stops = schedule.venuelist.venuesList.map(flattenStop);
      fetching = false;
    } catch (e) {
      if (!authHelper.isAuthorized($session)) {
        return goto(`${eventKey}/auth`);
      }
    }
  })

  // $: console.log($nextStop);
</script>

<style>
  .HOME {
    height: 90%;
  }
</style>

{#if fetching}
  <div class="flex h-full items-center text-center text-white text-6xl">
    Fetching Schedule
  </div>
{:else}
  <div class="HOME text-white text-center text-4xl pt-32">
    {#if $nextStop}
      <p>
        Time Remaining
      </p>
      <p class="text-orange-light">
        <Countdown to="{$nextStop.start}"/>
      </p>
      <p class="text-2xl mt-32">
        Next Bar: {$nextStop.name}<br>
        <a
          href="https://www.google.com/maps/place/{$nextStop.address}"
          target="_blank"
          class="text-blue-300 underline text-xl"
        >
          {$nextStop.address}
        </a>
      </p>
      <a
              href="{$page.path}/scores"
              class="block btn btn-primary w-2/3 my-16 mx-auto"
      >
        See scores
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
            <span class="text-extrabold uppercase">Current:</span>
          {/if}
          {stop.name}<br>
          <!--<a
            href="https://www.google.com/maps/place/{stop.address}"
            target="_blank"
            class="text-blue-400 underline"
          >-->
            {stop.address}
            <!--</a>-->
        </li>
      {/each}
    </ol>
  {/if}
{/if}

