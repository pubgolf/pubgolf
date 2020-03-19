const MS_IN_S = 1000;
const S_IN_M = 60;
const M_IN_H = 60;
const H_IN_D = 24;

/**
 * Get the time difference between two dates as an object keyed by units
 *
 * @param {Date} startDate - The start of the time interval
 * @param {Date} endDate - The end of the time interval
 *
 * @returns {{hours: number, seconds: number, minutes: number, days: number}}
 *          The delta between the two inputs, where each unit stacks
 *          on top of the bigger ones
 */
export function dateDelta (startDate, endDate) {
  const millisecondsTotal = endDate - startDate;
  const secondsTotal = Math.floor(millisecondsTotal / MS_IN_S);
  const minutesTotal = Math.floor(secondsTotal / S_IN_M);
  const hoursTotal = Math.floor(minutesTotal / M_IN_H);
  const daysTotal = Math.floor(hoursTotal / H_IN_D);

  return {
    days: daysTotal,
    hours: hoursTotal % H_IN_D,
    minutes: minutesTotal % M_IN_H,
    seconds: secondsTotal % S_IN_M,
  };
}

/**
 * Get a given number as a string, padded with zeroes if necessary
 *
 * @param {number} num - The number to format
 * @param {number} [minLength=2] - The shortest length of the output.
 *                                 Defaults to 2 if not provided
 *
 * @returns {string} The given number, prepended with zeroes
 *                   if it was fewer than `minLength` digits
 */
export function zeroPad (num, minLength = 2) {
  return `${num}`.padStart(minLength, '0');
}

/**
 * Determine from the session whether the server is in dev mode
 * @param {string} env
 *
 * @returns {boolean}
 */
export function isDev ({ config: { PUBGOLF_ENV: env = '' } }) {
  return env.endsWith('dev');
}

/**
 * Capitalize the first character and optionally lowercase the rest
 * @param {string} str
 * @param {boolean} [lowerRest]
 *
 * @returns {string}
 */
export function capFirst (str, lowerRest = false) {
  const rest = lowerRest
    ? str.slice(1).toLowerCase()
    : str.slice(1);
  return `${str[0].toUpperCase()}${rest}`;
}

/**
 * Naively capitalize a string
 * @param {string} str - Any old string
 *
 * @returns {string}
 */
export function capitalize (str) {
  return capFirst(str, true);
}

/**
 * Transform an object by applying a given function to each entry
 *
 * @template OriginalValueType, FinalValueType
 * @param {Object.<PropertyKey, OriginalValueType>} obj
 * @param {function([PropertyKey, OriginalValueType]): [PropertyKey, FinalValueType]} mapFn
 *
 * @returns {Object.<PropertyKey, FinalValueType>}
 */
export function mapEntries (obj, mapFn) {
  return Object.fromEntries(Object.entries(obj).map(mapFn));
}
