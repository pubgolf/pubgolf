import { describe, it, expect } from 'vitest';
import { formatListAnd, formatTimestamp, formatPlayerName } from './formatters';
import { ScoringCategory } from '$lib/proto/api/v1/shared_pb';
import type { Player } from '$lib/proto/api/v1/shared_pb';

function makePlayer(overrides: Partial<Player> = {}): Player {
	return {
		$typeName: 'api.v1.Player',
		id: 'player-1',
		data: {
			$typeName: 'api.v1.PlayerData',
			name: 'Alice',
			scoringCategory: ScoringCategory.UNSPECIFIED
		},
		events: [],
		...overrides
	} as Player;
}

describe('formatListAnd', () => {
	it('returns empty string for empty array', () => {
		expect(formatListAnd([])).toBe('');
	});

	it('returns the single item for a one-element array', () => {
		const result = formatListAnd(['one']);
		expect(result).toContain('one');
	});

	it('contains both items and "and" for two-element array', () => {
		const result = formatListAnd(['one', 'two']);
		expect(result).toContain('one');
		expect(result).toContain('two');
		expect(result).toContain('and');
	});

	it('contains all three items for a three-element array', () => {
		const result = formatListAnd(['a', 'b', 'c']);
		expect(result).toContain('a');
		expect(result).toContain('b');
		expect(result).toContain('c');
	});
});

describe('formatTimestamp', () => {
	it('returns a non-empty string for current date', () => {
		const result = formatTimestamp(new Date());
		expect(result).toBeTruthy();
		expect(typeof result).toBe('string');
	});

	it('does not throw for epoch date', () => {
		expect(() => formatTimestamp(new Date(0))).not.toThrow();
	});
});

describe('formatPlayerName', () => {
	it('includes scoring category display name for matching eventKey registration', () => {
		const player = makePlayer({
			events: [
				{
					$typeName: 'api.v1.EventRegistration',
					eventKey: 'event-2024',
					scoringCategory: ScoringCategory.PUB_GOLF_NINE_HOLE
				}
			]
		});
		const result = formatPlayerName(player, 'event-2024');
		expect(result).toContain('Alice');
		expect(result).toContain('Pro');
	});

	it('includes (None) when no matching registration', () => {
		const player = makePlayer();
		const result = formatPlayerName(player, 'event-2024');
		expect(result).toContain('Alice');
		expect(result).toContain('(None)');
	});
});
