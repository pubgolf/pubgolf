import { EVENT_KEY } from './mock-api';

export function makePlayer(overrides: Record<string, unknown> = {}) {
	return {
		id: 'player-1',
		data: { name: 'Alice Test' },
		events: [{ eventKey: EVENT_KEY, scoringCategory: 1 }],
		...overrides
	};
}

export function makeVenue(overrides: Record<string, unknown> = {}) {
	return {
		id: 'venue-1',
		name: 'The Anchor',
		address: '123 Main St',
		...overrides
	};
}

export function makeStage(overrides: Record<string, unknown> = {}) {
	return {
		id: 'stage-1',
		venue: makeVenue(),
		rule: { id: 'rule-1', venueDescription: 'Par 3 - drink in 3 sips' },
		rank: 1,
		durationMin: 45,
		...overrides
	};
}

export function makeStageScore(overrides: Record<string, unknown> = {}) {
	return {
		stageId: 'stage-1',
		playerId: 'player-1',
		score: { id: 'score-1', data: { value: 3 } },
		adjustments: [],
		isVerified: false,
		...overrides
	};
}

export function makeAdjustmentTemplate(overrides: Record<string, unknown> = {}) {
	return {
		id: 'template-1',
		data: {
			adjustment: { value: 1, label: 'Penalty: Wrong drink' },
			rank: 1,
			eventKey: EVENT_KEY,
			isVisible: true
		},
		...overrides
	};
}
