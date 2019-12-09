<script>
  const events = [
    { key: 'nyc-2019', text: 'NYC 2019' },
  ];
  if (process.env.NODE_ENV === 'development') {
    events.push(
      { key: 'sf-2019', text: 'SF 2019' },
      { key: 'current', text: 'Current' },
      { key: 'starts-in-30m', text: 'Starts in 30 Minutes' },
      { key: 'ended-30m-ago', text: 'Ended 30 Minutes Ago' },
    );
  }
</script>

<div class="flex justify-center items-center h-full">
  <div>
    <h1 class="text-4xl font-bold">
      Choose an Event
    </h1>

    <ul class="text-xl list-disc pl-8">
      {#each events as event (event.key)}
        <li class="mb-4">
          <a class="text-blue-400 underline" rel=prefetch href="{event.key}">
            {event.text}
          </a>
        </li>
      {/each}
    </ul>
  </div>
</div>
