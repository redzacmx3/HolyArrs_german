// Fork: German ships in the build. The fork adds its own lines, so the upstream
// German pack never fits this build's key list and would leave the app in English.
type Words = Record<string, string>;

const de: Words = Object.assign({}, ...Object.values(import.meta.glob<{ default: Words }>('./messages/de/*.json', { eager: true })).map((f) => f.default));

/** The built-in words of a language other than English (lines it lacks in English), or null when it comes as a pack. */
export function builtInWords(code: string, en: Words): Words | null {
	return code === 'de' ? { ...en, ...de } : null;
}
