import { describe, it, expect } from 'vitest';
import { create } from '@bufbuild/protobuf';
import { separateIds, combineIds } from './scores';
import { StageScoreSchema } from '$lib/proto/api/v1/admin_pb';

function makeStageScore() {
	return create(StageScoreSchema, {
		stageId: 'stage-1',
		playerId: 'player-1',
		score: {
			id: 'score-id-1',
			data: { value: 5 }
		},
		adjustments: [
			{ id: 'adj-id-1', data: { value: -1, label: 'Birdie' } },
			{ id: 'adj-id-2', data: { value: 1, label: 'Bogey' } }
		],
		isVerified: false
	});
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
		const s = create(StageScoreSchema, {
			stageId: 'stage-1',
			playerId: 'player-1',
			score: { id: 'score-id-1', data: { value: 5 } }
		});
		const result = separateIds(s);

		expect(result.data.adjustments).toHaveLength(0);
		expect(result.ids.adjustments).toHaveLength(0);
	});

	it('returns defaults for missing score', () => {
		const s = create(StageScoreSchema, {
			stageId: 'stage-1',
			playerId: 'player-1'
		});
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

	it('throws when ids.adjustments is shorter than data.adjustments', () => {
		const s = makeStageScore();
		const { data, ids } = separateIds(s);
		ids.adjustments.pop();
		expect(() => combineIds({ data, ids })).toThrow();
	});
});
