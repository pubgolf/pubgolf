import { describe, it, expect } from 'vitest';
import { scoringCategoryToDisplayName } from './scoring-category';
import { ScoringCategorySchema } from '$lib/proto/api/v1/shared_pb';

describe('scoringCategoryToDisplayName', () => {
	it('has an entry for every ScoringCategory enum value', () => {
		for (const enumValue of ScoringCategorySchema.values) {
			expect(
				scoringCategoryToDisplayName,
				`missing display name for ${enumValue.name} (${enumValue.number})`
			).toHaveProperty(String(enumValue.number));
		}
	});
});
