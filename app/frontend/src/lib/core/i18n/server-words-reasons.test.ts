// Fork: why a Release is turned down reads in German, and names the setting by its German UI name.
import { afterEach, describe, expect, it } from 'vitest';
import en from './messages/en/server-lines.json';
import de from './messages/de/server-lines.json';
import { lang } from './i18n.svelte.ts';
import { serverWords } from './server-words.ts';

const start = { value: lang.value, words: lang.words };
afterEach(() => {
	lang.value = start.value;
	lang.words = start.words;
});

describe('Release reasons in German', () => {
	it('reads the judge lines with their numbers', () => {
		lang.value = 'de';
		lang.words = { ...(en as Record<string, string>), ...(de as Record<string, string>) };
		for (const [line, want] of [
			['-23300 points, needs at least 11000', '-23300 Punkte, braucht mindestens 11000'],
			['no ladder step: the source is unknown', 'keine Stufe: die Quelle ist unbekannt'],
			['WEBDL-480p is not allowed in these Pick rules', 'WEBDL-480p ist in diesen Auswahlregeln nicht erlaubt'],
			['below good enough (Bluray-1080p)', 'unter „Gut genug“ (Bluray-1080p)'],
			['0.21 GB is under the 0.50 GB minimum', '0.21 GB liegt unter dem Minimum von 0.50 GB'],
			['S02E31 already has a copy as good', 'S02E31 hat schon eine gleich gute Kopie']
		])
			expect(serverWords(line)).toBe(want);
	});

	it('names the season pack setting by its German UI name and reads each episode of the list', () => {
		lang.value = 'de';
		lang.words = { ...(en as Record<string, string>), ...(de as Record<string, string>) };
		const line = 'not every episode would gain: S02E31 already has a copy as good; S02E32 already has a copy as good (Settings > Search > Season packs: A season pack is taken when every episode in it would be new or better)';
		expect(serverWords(line)).toBe(
			'nicht jede Folge würde besser: S02E31 hat schon eine gleich gute Kopie; S02E32 hat schon eine gleich gute Kopie (Einstellungen › Suche › Staffelpakete › Ein Staffelpaket wird genommen, wenn: Jede Folge darin neu oder besser wäre)'
		);
		expect(serverWords('not every episode would gain: S01E01 already has a copy as good')).toBe('nicht jede Folge würde besser: S01E01 hat schon eine gleich gute Kopie');
	});
});
