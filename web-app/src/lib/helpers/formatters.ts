import type { Player } from '$lib/proto/api/v1/shared_pb';
import { scoringCategoryToDisplayName } from './scoring-category';

const dateFormatter = new Intl.DateTimeFormat(undefined, { timeStyle: 'medium' });
export function formatTimestamp(d: Date) {
	return dateFormatter.format(d);
}

const listFormatter = new Intl.ListFormat(undefined, { type: 'conjunction' });
export function formatListAnd(l: string[]) {
	return listFormatter.format(l);
}
export function formatPlayerName(player: Player, eventKey: string) {
	const cat = player?.events.find((x) => x.eventKey == eventKey)?.scoringCategory;
	const catName = cat ? scoringCategoryToDisplayName[cat] : 'None';
	return `${player?.data?.name} (${catName})`;
}
