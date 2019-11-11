import { readable, writable } from 'svelte/store';

/**
 * Global stores that can be subscribed to from anywhere in the app
 */


/**
 * The current time
 * @type {Readable<Date>}
 */
export const time = readable(new Date(), function start(set) {
  const interval = setInterval(() => {
    set(new Date());
  }, 1000);

  return function stop() {
    clearInterval(interval);
  };
});
