import { describe, it, expect } from 'vitest';
import { separateIds, combineIds } from './scores';
import type { StageScore } from '$lib/proto/api/v1/admin_pb';

// Construct a minimal StageScore fixture that matches the proto-generated shape.
// We use plain objects cast to the type rather than proto create() to avoid
// runtime proto registry setup in unit tests.

function makeStageScore(overrides: Partial<StageScore> = {}): StageScore {
	return {
		$typeName: 'api.v1.StageScore',
		stageId: 'stage-1',
		playerId: 'player-1',
		score: {
			$typeName: 'api.v1.Score',
			id: 'score-id-1',
			data: {
				$typeName: 'api.v1.ScoreData',
				value: 5
			}
		},
		adjustments: [
			{
				$typeName: 'api.v1.Adjustment',
				id: 'adj-id-1',
				data: {
					$typeName: 'api.v1.AdjustmentData',
					value: -1,
					label: 'Birdie'
				}
			},
			{
				$typeName: 'api.v1.Adjustment',
				id: 'adj-id-2',
				data: {
					$typeName: 'api.v1.AdjustmentData',
					value: 1,
					label: 'Bogey'
				}
			}
		],
		isVerified: false,
		...overrides
	} as StageScore;
}

describe('separateIds', () => {
	it('separates ids and data from a full StageScore', () => {
		const s = makeStageScore();
		const result = separateIds(s);

		// data side should carry values but not ids
		expect(result.data.stageId).toBe('stage-1');
		expect(result.data.playerId).toBe('player-1');
		expect(result.data.score?.value).toBe(5);
		expect(result.data.adjustments).toHaveLength(2);
		expect(result.data.adjustments?.[0]).toMatchObject({ value: -1, label: 'Birdie' });
		expect(result.data.adjustments?.[1]).toMatchObject({ value: 1, label: 'Bogey' });

		// ids side should carry ids
		expect(result.ids.score.id).toBe('score-id-1');
		expect(result.ids.adjustments).toHaveLength(2);
		expect(result.ids.adjustments[0].id).toBe('adj-id-1');
		expect(result.ids.adjustments[1].id).toBe('adj-id-2');
	});

	it('returns empty adjustment arrays when there are no adjustments', () => {
		const s = makeStageScore({ adjustments: [] });
		const result = separateIds(s);

		expect(result.data.adjustments).toHaveLength(0);
		expect(result.ids.adjustments).toHaveLength(0);
	});

	it('returns defaults for missing score', () => {
		const s = makeStageScore({ score: undefined });
		const result = separateIds(s);

		expect(result.data.score?.value).toBe(0);
		expect(result.ids.score.id).toBe('');
	});

	it('round-trips through combineIds', () => {
		const s = makeStageScore();
		const { data, ids } = separateIds(s);
		const combined = combineIds({ data, ids });

		expect(combined.stageId).toBe(s.stageId);
		expect(combined.playerId).toBe(s.playerId);
		expect(combined.score?.id).toBe(s.score?.id);
		expect(combined.score?.data?.value).toBe(s.score?.data?.value);
		expect(combined.adjustments).toHaveLength(2);
		expect(combined.adjustments[0].id).toBe('adj-id-1');
		expect(combined.adjustments[1].id).toBe('adj-id-2');
	});
});

describe('combineIds', () => {
	it('aligns ids and data arrays by index for multiple adjustments', () => {
		const s = makeStageScore();
		const { data, ids } = separateIds(s);
		const combined = combineIds({ data, ids });

		expect(combined.adjustments[0].data).toMatchObject({ value: -1, label: 'Birdie' });
		expect(combined.adjustments[0].id).toBe('adj-id-1');
		expect(combined.adjustments[1].data).toMatchObject({ value: 1, label: 'Bogey' });
		expect(combined.adjustments[1].id).toBe('adj-id-2');
	});
});
