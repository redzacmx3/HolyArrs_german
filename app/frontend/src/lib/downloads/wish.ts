// Language and quality wishes: a small change laid over the Pick rules for
// every movie or show, a title, a season or an episode.

/** A wish as the server keeps it; a missing field keeps the level above. */
export type Wish = { languages?: string[] | null; original_ok?: boolean | null; quality?: string };

/** What a set of Pick rules asks for, in short. */
export type WishRules = { name: string; allowed: string[]; good_enough_at: string; languages: string[]; original_ok: boolean };

/**
 * The language choices the assistant offers: keep the rules' own, a language
 * or the original, only a language, or only the original language.
 */
export type LangChoice = 'keep' | 'with_original' | 'only' | 'original';

/** The wish for the choices made; quality '' keeps the rules' steps. */
export function wishBody(lang: LangChoice, code: string, quality: string): Wish {
	const w: Wish = {};
	if (lang === 'with_original') Object.assign(w, { languages: [code], original_ok: true });
	if (lang === 'only') Object.assign(w, { languages: [code], original_ok: false });
	if (lang === 'original') Object.assign(w, { languages: [], original_ok: true });
	if (quality) w.quality = quality;
	return w;
}

/** The choices that make up a saved wish (none saved: keep everything). */
export function wishFrom(w: Wish | undefined | null): { lang: LangChoice; code: string; quality: string } {
	const quality = w?.quality ?? '';
	const langs = w?.languages;
	if (!w || langs == null) return { lang: 'keep', code: '', quality };
	if (langs.length === 0) return { lang: 'original', code: '', quality };
	return { lang: w.original_ok === false ? 'only' : 'with_original', code: langs[0], quality };
}
