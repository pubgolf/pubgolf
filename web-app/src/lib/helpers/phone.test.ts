import { describe, it, expect } from 'vitest';
import { cleanPhoneNumber } from './phone';

describe('cleanPhoneNumber', () => {
	it('strips dashes and adds +1 prefix', () => {
		expect(cleanPhoneNumber('555-867-5309')).toBe('+15558675309');
	});

	it('strips parentheses and spaces and adds +1 prefix', () => {
		expect(cleanPhoneNumber('(555) 867-5309')).toBe('+15558675309');
	});

	it('prepends + when 11 digits starting with 1', () => {
		expect(cleanPhoneNumber('15558675309')).toBe('+15558675309');
	});

	it('adds 1 for 10-digit number', () => {
		expect(cleanPhoneNumber('5558675309')).toBe('+15558675309');
	});

	it('handles empty input (documents current behavior)', () => {
		expect(cleanPhoneNumber('')).toBe('+1');
	});
});
