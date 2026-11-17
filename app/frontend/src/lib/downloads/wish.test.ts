import { describe, expect, it } from 'vitest';
import { wishBody, wishFrom } from './wish.ts';

describe('wishes', () => {
	it('builds the wish for each language choice', () => {
		expect(wishBody('keep', 'de', '')).toEqual({});
		expect(wishBody('with_original', 'de', '1080p')).toEqual({ languages: ['de'], original_ok: true, quality: '1080p' });
		expect(wishBody('only', 'en', '')).toEqual({ languages: ['en'], original_ok: false });
		expect(wishBody('original', 'de', '2160p')).toEqual({ languages: [], original_ok: true, quality: '2160p' });
	});
	it('reads a saved wish back into the same choices', () => {
		for (const [lang, code, q] of [['with_original', 'de', '720p'], ['only', 'en', ''], ['original', '', '1080p'], ['keep', '', '2160p']] as const) {
			expect(wishFrom(wishBody(lang, code, q))).toEqual({ lang, code, quality: q });
		}
		expect(wishFrom(undefined)).toEqual({ lang: 'keep', code: '', quality: '' });
		expect(wishFrom({ languages: null, quality: '' })).toEqual({ lang: 'keep', code: '', quality: '' });
	});
});
