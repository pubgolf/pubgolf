// What index of the cleaned value to insert a string before
const insertions = { 0: '(', 3: ') ', 6: '-' };
// The maximum length WITHOUT insertions
const maxPrunedLength = 10;

/**
 * A map function for a strings
 * @param {string} original The string to map on
 * @param {function} fn The function to apply to each character of the string
 * @returns {string}
 */
function mapString (original, fn) {
  return original.split('').map(fn).join('');
}

/**
 * Strip all non-numeric characters out of a string
 * @param {string} value
 * @returns {string} The string as purely numeric
 */
export function onlyDigits (value = '') {
  if (typeof value === 'string') {
    const pruned = value.replace(/\D/g, '');
    // Remove leading 1
    const start = pruned.startsWith('1') ? 1 : 0;
    return pruned.slice(start, maxPrunedLength + start);
  }
  return '';
}

/**
 * Apply the insertions defined above
 * @param {string} value - The original value
 * @returns {string} The string with the insertions added
 */
export function applyInsertions (value = '') {
  const digits = onlyDigits(value);
  if (digits) {
    return mapString(digits, (character, index) => {
      const insertion = insertions[index] || '';
      return insertion + character;
    });
  }
  return '';
}
