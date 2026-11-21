// The number release names use for one episode, shown on the episode row so a
// wrong match can be spotted. Release names follow the show's
// release_numbers setting: '' the order on screen, 'tmdb' TMDB's own numbers,
// 'tvmaze' TVmaze's (like TVDB).
import { t } from '#lib/core/i18n/i18n.svelte.ts';

export interface NumberedEpisode {
	season?: number;
	episode: number;
	tmdb_season?: number;
	tmdb_episode?: number;
	tvmaze_season?: number;
	tvmaze_episode?: number;
}

const code = (s: number, e: number) => `S${String(s).padStart(2, '0')}E${String(e).padStart(2, '0')}`;

/** The release-name code ("S02E41"), 'missing' when TVmaze has no match, or null when names follow the numbers on screen. */
export function releaseCode(e: NumberedEpisode, mode: string | undefined, season = e.season ?? 0): string | null {
	if (!mode || season === 0) return null;
	if (mode === 'tmdb') return code(e.tmdb_season ?? season, e.tmdb_episode ?? e.episode);
	if (mode === 'tvmaze') return e.tvmaze_season == null || e.tvmaze_episode == null ? 'missing' : code(e.tvmaze_season, e.tvmaze_episode);
	return null;
}

/** The words for the episode row, or '' when there is nothing to show. */
export function releaseNumberLine(e: NumberedEpisode, mode: string | undefined, season?: number): string {
	const c = releaseCode(e, mode, season ?? e.season ?? 0);
	if (!c) return '';
	return c === 'missing' ? t('discover.titleLibrary.releaseNumberMissing') : t('discover.titleLibrary.releaseNumberIs', { code: c });
}
