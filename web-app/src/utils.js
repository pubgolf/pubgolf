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

export function isDev () {
  return process.env.NODE_ENV === 'development';
}
