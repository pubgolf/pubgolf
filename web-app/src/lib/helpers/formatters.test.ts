import { describe, it, expect } from 'vitest';
import { create } from '@bufbuild/protobuf';
import { formatPlayerName } from './formatters';
import { PlayerSchema, ScoringCategory } from '$lib/proto/api/v1/shared_pb';

describe('formatPlayerName', () => {
	it('includes scoring category display name for matching eventKey registration', () => {
		const player = create(PlayerSchema, {
			id: 'player-1',
			data: { name: 'Alice' },
			events: [
				{
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
		const player = create(PlayerSchema, {
			id: 'player-1',
			data: { name: 'Alice' }
		});
		const result = formatPlayerName(player, 'event-2024');
		expect(result).toContain('Alice');
		expect(result).toContain('(None)');
	});
});
