<script>
  import { time } from '../stores';
  import { dateDelta, zeroPad } from '../utils';

  // `export let` is used to declare component props
  export let to; // What date to count down to

  // `$time` is an auto-subscribed reference to the value of the `time` store
  // Code after `$:` is re-run anytime one of its dependencies changes.
  // So these are reactive assignments (like computed properties in Vue)
  $: delta = dateDelta($time, to);
</script>

{#if delta.hours > 0}
  {delta.days}d {delta.hours}h {delta.minutes}m
{:else}
  {delta.minutes}:{zeroPad(delta.seconds)}
{/if}
