import { describe, it, expect } from 'vitest';
import { cleanPhoneNumber } from './phone';

describe('cleanPhoneNumber', () => {
	it.each([
		{ input: '555-867-5309', expected: '+15558675309', label: 'strips dashes and adds +1 prefix' },
		{
			input: '(555) 867-5309',
			expected: '+15558675309',
			label: 'strips parentheses and spaces'
		},
		{
			input: '15558675309',
			expected: '+15558675309',
			label: 'prepends + when 11 digits starting with 1'
		},
		{ input: '5558675309', expected: '+15558675309', label: 'adds 1 for 10-digit number' },
		{ input: '', expected: '', label: 'returns empty for empty input' }
	])('$label: "$input" -> "$expected"', ({ input, expected }) => {
		expect(cleanPhoneNumber(input)).toBe(expected);
	});
});
