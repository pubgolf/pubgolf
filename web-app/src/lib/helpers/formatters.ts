const dateFormatter = new Intl.DateTimeFormat(undefined, { timeStyle: 'medium' });
export function formatTimestamp(d: Date) {
	return dateFormatter.format(d);
}

const listFormatter = new Intl.ListFormat(undefined, { type: 'conjunction' });
export function formatListAnd(l: string[]) {
	return listFormatter.format(l);
}
