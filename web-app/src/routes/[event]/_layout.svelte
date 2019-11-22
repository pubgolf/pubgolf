<script>
  import { stores } from '@sapper/app';
  import { onMount } from 'svelte';

  import {
    api,
    event,
  } from '../../stores';
  import { isDev } from '../../utils';

  const { page, session } = stores();

  onMount(() => {
    const { params } = $page;

    $event = params.event;

    if (isDev()) {
      window.$api = $api;
    }
  });
</script>

<slot></slot>
