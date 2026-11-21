import { describe, expect, it } from 'vitest';
import { releaseCode } from './release-number.ts';

describe('releaseCode', () => {
	const eagle = { season: 3, episode: 47, tmdb_season: 2, tmdb_episode: 45, tvmaze_season: 2, tvmaze_episode: 41 };
	it('shows nothing when release names follow the order on screen', () => {
		expect(releaseCode(eagle, '')).toBeNull();
		expect(releaseCode(eagle, undefined)).toBeNull();
	});
	it('reads TMDB numbers, or the numbers on screen when TMDB numbers the same', () => {
		expect(releaseCode(eagle, 'tmdb')).toBe('S02E45');
		expect(releaseCode({ season: 1, episode: 3 }, 'tmdb')).toBe('S01E03');
	});
	it('reads TVmaze numbers and says when TVmaze has no match', () => {
		expect(releaseCode(eagle, 'tvmaze')).toBe('S02E41');
		expect(releaseCode({ season: 7, episode: 2 }, 'tvmaze')).toBe('missing');
	});
	it('takes the season from the list when the episode has none, and leaves specials out', () => {
		expect(releaseCode({ episode: 4, tvmaze_season: 5, tvmaze_episode: 9 }, 'tvmaze', 6)).toBe('S05E09');
		expect(releaseCode({ season: 0, episode: 1 }, 'tvmaze')).toBeNull();
	});
});
