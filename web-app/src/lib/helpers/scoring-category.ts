import { ScoringCategory } from '$lib/proto/api/v1/shared_pb';

export const scoringCategoryToDisplayName = {
	[ScoringCategory.UNSPECIFIED]: '',
	[ScoringCategory.PUB_GOLF_NINE_HOLE]: 'Pro',
	[ScoringCategory.PUB_GOLF_FIVE_HOLE]: 'Semi-Pro',
	[ScoringCategory.PUB_GOLF_CHALLENGES]: 'Masters'
};
