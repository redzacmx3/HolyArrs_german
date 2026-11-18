import { describe, expect, it } from 'vitest';
import { fits, noPicks, options, valuesOf } from './release-filters.ts';
import type { ReleaseTags } from '../core/types.ts';

const r = (tags?: ReleaseTags) => ({ tags });
const list = [
	r({ languages: ['de'], dual: true, resolution: 1080, source: 'web_dl', codec: 'h264', audio: 'eac3', channels: '5.1' }),
	r({ languages: ['de'], resolution: 720, source: 'web_dl', codec: 'h264', audio: 'ac3' }),
	r({ resolution: 2160, source: 'web_dl', codec: 'h265', hdr: ['dv', 'hdr10'], audio: 'eac3', atmos: true }),
	r()
];

describe('release filters', () => {
	it('reads each facet from the name tags, and none when the name says nothing', () => {
		expect(valuesOf(list[0], 'lang')).toEqual(['de', 'dl']);
		expect(valuesOf(list[2], 'lang')).toEqual(['none']);
		expect(valuesOf(list[2], 'hdr')).toEqual(['dv', 'hdr10']);
		expect(valuesOf(list[1], 'hdr')).toEqual(['sdr']);
		expect(valuesOf(list[2], 'audio')).toEqual(['eac3', 'atmos']);
		expect(valuesOf(list[3], 'res')).toEqual(['none']);
	});
	it('combines the picks of every facet', () => {
		expect(list.filter((x) => fits(x, { ...noPicks(), lang: 'de' }))).toHaveLength(2);
		expect(list.filter((x) => fits(x, { ...noPicks(), lang: 'de', res: '1080' }))).toEqual([list[0]]);
		expect(list.filter((x) => fits(x, { ...noPicks(), audio: 'atmos' }))).toEqual([list[2]]);
	});
	it('counts each option under the other picks, German first and unknown last, and keeps a picked value with none left', () => {
		const lang = options(list, 'lang', noPicks());
		expect(lang.map((o) => [o.value, o.n])).toEqual([['de', 2], ['dl', 1], ['none', 2]]);
		const res = options(list, 'res', { ...noPicks(), lang: 'de' });
		expect(res.map((o) => [o.value, o.n])).toEqual([['2160', 0], ['1080', 1], ['720', 1], ['none', 0]]);
		expect(options(list, 'codec', { ...noPicks(), codec: 'av1' }).some((o) => o.value === 'av1' && o.n === 0)).toBe(true);
	});
});
