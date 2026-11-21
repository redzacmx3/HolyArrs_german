// Fork: what each trait added to or took from a Release's points, for the
// rows in "Release auswählen". The costliest minus comes first.
import type { Reason } from '../core/types.ts';
import { appLocale } from '#lib/core/i18n/i18n.svelte.ts';

export interface PointsLine {
	name: string;
	points: number;
}

/** The trait lines of a Release's reasons, minus first (largest first), then plus (largest first). */
export function pointsLines(reasons: Reason[]): PointsLine[] {
	return reasons
		.filter((x) => (x.effect === 'add' || x.effect === 'subtract') && x.points)
		.map((x) => ({ name: x.text.replace(/ \([+-]\d+ points\)$/, '').replace(/: fits$/, ''), points: x.points! }))
		.sort((a, b) => Number(a.points > 0) - Number(b.points > 0) || Math.abs(b.points) - Math.abs(a.points));
}

/** A signed number: "+11,000", "−35,000". */
export function signed(n: number, locale = appLocale()): string {
	const s = Math.abs(n).toLocaleString(locale);
	return n < 0 ? `−${s}` : `+${s}`;
}
