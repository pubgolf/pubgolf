<script>
  import { stores } from '@sapper/app';
  import { onMount } from 'svelte';

  import {
    api,
    event,
  } from '../../stores';
  import { isDev } from '../../utils';
  import Nav from '../../components/Nav.svelte';


  export let segment;

  const { page, session } = stores();
  const links = [
    {
      segment: 'home',
      text: 'Home',
      icon: require('./_icons/beer.svg'),
    },
    // {
    //   segment: 'leaderboard',
    //   text: 'Scores',
    //   icon: require('./_icons/leaderboard.svg'),
    // },
    // {
    //   segment: 'admin',
    //   text: 'Admin',
    //   icon: require('./_icons/clipboard.svg'),
    // },
  ];

  onMount(() => {
    const { params } = $page;

    $event = params.event;

    if (isDev()) {
      window.$api = $api;
    }
  });
</script>

<style>
  .EVENT-LAYOUT {
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  .MAIN-WRAPPER {
    flex: 1 1 auto;
    overflow-y: scroll;
  }

  .NAV-WRAPPER {
    flex: 0 0 auto;
  }
</style>

{#if segment === 'auth'}
  <slot></slot>
{:else}
  <div class="EVENT-LAYOUT">
    <div class="MAIN-WRAPPER">
      <slot></slot>
    </div>
    <div class="NAV-WRAPPER">
      <Nav basePath="{$event}" {links} {segment}/>
    </div>
  </div>
{/if}
