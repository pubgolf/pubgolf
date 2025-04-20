import { StageScoreDataSchema, type StageScore } from '$lib/proto/api/v1/admin_pb';
import type { MessageInitShape } from '@bufbuild/protobuf';

export type StageScoreIds = {
	score: {
		id: string;
	};
	adjustments: {
		id: string;
	}[];
};

export function separateIds(s: StageScore): {
	data: MessageInitShape<typeof StageScoreDataSchema>;
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
	data: MessageInitShape<typeof StageScoreDataSchema>;
	ids: StageScoreIds;
}) {
	return {
		stageId: data.stageId,
		playerId: data.playerId,
		score: {
			id: ids.score.id,
			data: data.score
		},
		adjustments: Array.from(Array(data.adjustments?.length || 0).keys()).map((i) => ({
			id: ids.adjustments[i].id,
			data: data.adjustments?.[i]
		}))
	};
}
