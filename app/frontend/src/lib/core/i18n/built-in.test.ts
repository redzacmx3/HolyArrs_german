import { afterEach, describe, expect, it, vi } from 'vitest';
import { lang, langPack, setLang } from './i18n.svelte.ts';

// Fork: German ships in the build, so picking it never asks the server for a pack.
describe('built-in German', () => {
	afterEach(() => vi.unstubAllGlobals());
	it('shows German without fetching a pack', async () => {
		const fetch = vi.fn(async () => new Response('{}', { status: 404 }));
		vi.stubGlobal('fetch', fetch);
		vi.stubGlobal('localStorage', { getItem: () => null, setItem: () => {}, removeItem: () => {} });
		await setLang('de');
		expect(fetch).not.toHaveBeenCalled();
		expect(lang.value).toBe('de');
		expect(langPack.missing).toBeNull();
		expect(lang.words['core.format.justNow']).toBe('gerade eben');
		await setLang('en');
	});
});
