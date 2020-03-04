import {
  derived,
  readable,
  writable,
} from 'svelte/store';

/**
 * Global stores that can be subscribed to from anywhere in the app
 */


/**
 * The current time
 * @type {Readable<Date>}
 */
export const time = readable(new Date(), function start (set) {
  const interval = setInterval(() => {
    set(new Date());
  }, 1000);

  return function stop () {
    clearInterval(interval);
  };
});

export const event = writable('');
export const api = writable(null);

export const stops = writable([]);
export const nextStop = derived(
  [
    time,
    stops,
  ],
  // TODO: Use a lazier timer
  ([$time, $stops]) => (
    ($stops && $stops.length
    )
      ? $stops.find(stop => stop.start > $time)
      : null
  ),
);
export const pastStops = derived(
  [
    stops,
    nextStop,
  ],
  ([$stops, $nextStop]) => (
    $nextStop ? $stops.slice(0, $nextStop.index).reverse() : $stops.reverse()
  ),
);
