import {
  derived,
  readable,
  writable,
} from 'svelte/store';

import { getAPI } from './api';

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
export const api = derived(
  event,
  $event => getAPI($event),
);
