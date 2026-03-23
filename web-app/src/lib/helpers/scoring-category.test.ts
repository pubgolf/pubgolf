import { describe, it, expect } from 'vitest';
import { scoringCategoryToDisplayName } from './scoring-category';
import { ScoringCategory } from '$lib/proto/api/v1/shared_pb';

describe('scoringCategoryToDisplayName', () => {
	it('has an entry for every ScoringCategory enum value', () => {
		expect(scoringCategoryToDisplayName).toHaveProperty(String(ScoringCategory.UNSPECIFIED));
		expect(scoringCategoryToDisplayName).toHaveProperty(String(ScoringCategory.PUB_GOLF_NINE_HOLE));
		expect(scoringCategoryToDisplayName).toHaveProperty(String(ScoringCategory.PUB_GOLF_FIVE_HOLE));
		expect(scoringCategoryToDisplayName).toHaveProperty(
			String(ScoringCategory.PUB_GOLF_CHALLENGES)
		);
	});

	it('maps UNSPECIFIED to empty string', () => {
		expect(scoringCategoryToDisplayName[ScoringCategory.UNSPECIFIED]).toBe('');
	});

	it('maps PUB_GOLF_NINE_HOLE to "Pro"', () => {
		expect(scoringCategoryToDisplayName[ScoringCategory.PUB_GOLF_NINE_HOLE]).toBe('Pro');
	});

	it('maps PUB_GOLF_FIVE_HOLE to "Semi-Pro"', () => {
		expect(scoringCategoryToDisplayName[ScoringCategory.PUB_GOLF_FIVE_HOLE]).toBe('Semi-Pro');
	});

	it('maps PUB_GOLF_CHALLENGES to "Masters"', () => {
		expect(scoringCategoryToDisplayName[ScoringCategory.PUB_GOLF_CHALLENGES]).toBe('Masters');
	});
});
