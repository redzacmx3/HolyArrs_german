// Drop-down filters for "Release auswählen": language, resolution, source, HDR, codec and sound, read from what each Release name says (ReleaseView.tags).
import type { ReleaseView } from '../core/types.ts';
import { t } from '../core/i18n/i18n.svelte.ts';
import { languageName } from './copies.ts';

export type FacetKey = 'lang' | 'res' | 'source' | 'hdr' | 'codec' | 'audio';
export const FACETS: FacetKey[] = ['lang', 'res', 'source', 'hdr', 'codec', 'audio'];
/** The picked value per facet; '' is every Release. */
export type Picks = Record<FacetKey, string>;
export const noPicks = (): Picks => ({ lang: '', res: '', source: '', hdr: '', codec: '', audio: '' });

// A name that says nothing for a facet counts as "none" there, so it can be picked or left out too.
const NONE = 'none';

/** The values a Release has for a facet; a Release can have several (German and English, Atmos besides DD+). */
export function valuesOf(r: Pick<ReleaseView, 'tags'>, f: FacetKey): string[] {
	const g = r.tags ?? {};
	switch (f) {
		case 'lang': {
			const out = [...(g.languages ?? [])];
			if (g.dual) out.push('dl');
			if (g.subs?.length) out.push(...g.subs.map((s) => `sub:${s}`));
			return out.length ? out : [NONE];
		}
		case 'res':
			return [g.resolution ? String(g.resolution) : NONE];
		case 'source':
			return [g.source || NONE];
		case 'hdr':
			return g.hdr?.length ? g.hdr : ['sdr'];
		case 'codec':
			return [g.codec || NONE];
		case 'audio': {
			const out = g.audio ? [g.audio] : [];
			if (g.atmos) out.push('atmos');
			return out.length ? out : [NONE];
		}
	}
}

/** Whether a Release fits every picked value. */
export function fits(r: Pick<ReleaseView, 'tags'>, picks: Picks): boolean {
	return FACETS.every((f) => !picks[f] || valuesOf(r, f).includes(picks[f]));
}

export const picked = (picks: Picks) => FACETS.filter((f) => picks[f]).length;

export type Option = { value: string; label: string; n: number };

// Order inside a facet: best picture first, German first, unknown last.
const ORDER: Partial<Record<FacetKey, string[]>> = {
	lang: ['de', 'dl', 'en'],
	res: ['2160', '1080', '720', '576', '480'],
	source: ['remux', 'bluray', 'web_dl', 'web_rip', 'tv', 'dvd', 'disc_image'],
	hdr: ['dv', 'hdr10plus', 'hdr10', 'hdr', 'hlg', 'sdr'],
	codec: ['h265', 'h264', 'av1'],
	audio: ['atmos', 'truehd', 'dts_x', 'dts_hd_ma', 'dts', 'eac3', 'ac3', 'flac', 'aac']
};

/** The options of one facet: each value the list holds, with how many Releases have it once the other facets' picks apply. A picked value stays even when no Release has it. */
export function options(list: Pick<ReleaseView, 'tags'>[], f: FacetKey, picks: Picks): Option[] {
	const others = { ...picks, [f]: '' };
	const n = new Map<string, number>();
	for (const r of list) {
		const vs = valuesOf(r, f);
		const ok = fits(r, others);
		for (const v of vs) n.set(v, (n.get(v) ?? 0) + (ok ? 1 : 0));
	}
	if (picks[f] && !n.has(picks[f])) n.set(picks[f], 0);
	const rank = (v: string) => {
		if (v === NONE) return 1000;
		const i = ORDER[f]?.indexOf(v.replace(/^sub:/, '')) ?? -1;
		return (v.startsWith('sub:') ? 500 : 0) + (i < 0 ? 400 : i);
	};
	return [...n.entries()].map(([value, c]) => ({ value, label: valueLabel(f, value), n: c })).sort((a, b) => rank(a.value) - rank(b.value) || a.label.localeCompare(b.label));
}

const SOURCE: Record<string, string> = { web_dl: 'WEB-DL', web_rip: 'WEBRip', bluray: 'BluRay', remux: 'Remux', tv: 'TV', dvd: 'DVD', disc_image: 'Disc', cam: 'CAM', telesync: 'Telesync' };
const HDR: Record<string, string> = { dv: 'Dolby Vision', hdr10plus: 'HDR10+', hdr10: 'HDR10', hdr: 'HDR', hlg: 'HLG', sdr: 'SDR' };
const CODEC: Record<string, string> = { h264: 'x264 / H.264', h265: 'x265 / HEVC', av1: 'AV1', xvid: 'XviD', mpeg2: 'MPEG-2', vc1: 'VC-1' };
const AUDIO: Record<string, string> = { aac: 'AAC', ac3: 'DD (AC3)', eac3: 'DD+ (EAC3)', dts: 'DTS', dts_hd_ma: 'DTS-HD MA', dts_x: 'DTS:X', truehd: 'TrueHD', flac: 'FLAC', mp3: 'MP3', opus: 'Opus', pcm: 'PCM', atmos: 'Atmos' };

/** The words for one value in a facet's list. */
export function valueLabel(f: FacetKey, v: string): string {
	if (v === NONE) return f === 'lang' ? t('mediafiles.releaseFilter.noLanguage') : t('mediafiles.releaseFilter.unknown');
	switch (f) {
		case 'lang':
			if (v === 'dl') return t('mediafiles.releaseFilter.dual');
			if (v.startsWith('sub:')) return t('mediafiles.releaseFilter.subs', { language: languageName(v.slice(4)) });
			return languageName(v);
		case 'res':
			return v === '2160' ? '4K (2160p)' : `${v}p`;
		case 'source':
			return SOURCE[v] ?? v;
		case 'hdr':
			return HDR[v] ?? v.toUpperCase();
		case 'codec':
			return CODEC[v] ?? v.toUpperCase();
		case 'audio':
			return AUDIO[v] ?? v.toUpperCase();
	}
}

/** The facet's name above its list. */
export function facetLabel(f: FacetKey): string {
	return { lang: t('mediafiles.releaseFilter.lang'), res: t('mediafiles.releaseFilter.res'), source: t('mediafiles.releaseFilter.source'), hdr: t('mediafiles.releaseFilter.hdr'), codec: t('mediafiles.releaseFilter.codec'), audio: t('mediafiles.releaseFilter.audio') }[f];
}

// The picks stay for the next sheet in this browser, so "German" need not be picked for every episode.
const KEY = 'holyarrs.releaseFilters';
export function loadPicks(): Picks {
	try {
		const v = JSON.parse(localStorage.getItem(KEY) ?? '{}') as Partial<Picks>;
		const p = noPicks();
		for (const f of FACETS) if (typeof v[f] === 'string') p[f] = v[f];
		return p;
	} catch {
		return noPicks();
	}
}
export function savePicks(p: Picks): void {
	try {
		localStorage.setItem(KEY, JSON.stringify(p));
	} catch {
		// Private windows may refuse; the picks then last for this sheet only.
	}
}
