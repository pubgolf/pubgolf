import type { StageScore, StageScoreData } from '$lib/proto/api/v1/admin_pb';
import type { PartialMessage } from '@bufbuild/protobuf';
import type { Strict } from './types';

export type StageScoreIds = {
	score: {
		id: string;
	};
	adjustments: {
		id: string;
	}[];
};

export function separateIds(s: StageScore): {
	data: Strict<StageScoreData>;
	ids: StageScoreIds;
} {
	return {
		data: {
			stageId: s.stageId,
			playerId: s.playerId,
			score: {
				value: s.score?.data?.value || 0
			},
			adjustments: [...s.adjustments.map((a) => ({ value: 0, label: '', ...a.data }))]
		},
		ids: {
			score: {
				id: s.score?.id || ''
			},
			adjustments: [...s.adjustments.map((a) => ({ id: a.id || '' }))]
		}
	};
}

export function combineIds({
	data,
	ids
}: {
	data: Strict<StageScoreData>;
	ids: StageScoreIds;
}): PartialMessage<StageScore> {
	return {
		stageId: data.stageId,
		playerId: data.playerId,
		score: {
			id: ids.score.id,
			data: data.score
		},
		adjustments: Array.from(Array(data.adjustments.length).keys()).map((i) => ({
			id: ids.adjustments[i].id,
			data: data.adjustments[i]
		}))
	};
}
