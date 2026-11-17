package releasename

import "slices"

// subWords after a language word say the language is the subtitles'.
var subWords = setOf("sub", "subs", "subbed", "subtitled")

// readLanguages reads the audio and subtitle languages from the lower-case
// tags, and sets the multi, dual and omu flags.
func readLanguages(r *Result, low []string, flags map[string]bool, has func(...string) bool) {
	// Languages. A newer language word before an ordinary word belongs to
	// an episode title or a subtitle note ("The Chinese Restaurant", "Thai Sub").
	// A language word before a subtitle word names subtitles, not audio
	// ("German.Subbed", "GERMAN.SUBS").
	for i, t := range low {
		code, ok := langWords[t]
		if !ok {
			code, ok = moreLangWords[t]
			ok = ok && (i+1 == len(low) || !plainWord(low[i+1]))
			// "hun" also needs a tag before it, so "The.Hun.720p" stays a title.
			ok = ok && (t != "hun" || i > 0 && !plainWord(low[i-1]))
		}
		if !ok {
			continue
		}
		if i+1 < len(low) && subWords[low[i+1]] {
			if !slices.Contains(r.SubLanguages, code) {
				r.SubLanguages = append(r.SubLanguages, code)
			}
			continue
		}
		if !slices.Contains(r.Languages, code) {
			r.Languages = append(r.Languages, code)
		}
	}
	// OmU (Original mit Untertiteln): the original audio with German subtitles.
	if has("omu") {
		flags["omu"] = true
		r.Languages = slices.DeleteFunc(r.Languages, func(l string) bool { return l == "de" })
		if !slices.Contains(r.SubLanguages, "de") {
			r.SubLanguages = append(r.SubLanguages, "de")
		}
	}
	if has("multi") {
		flags["multi"] = true
	}
	// DL is dual language, except after WEB, where it is the source ("WEB.DL").
	for i, t := range low {
		if t == "dual" || t == "dl" && (i == 0 || low[i-1] != "web") {
			flags["dual"] = true
		}
	}
	// German DL holds German and the original audio; German ML holds German,
	// the original and English (the old apps' LanguageParser reads them so).
	if slices.Equal(r.Languages, []string{"de"}) {
		switch {
		case flags["dual"]:
			r.OriginalAudio = true
		case has("ml"):
			r.OriginalAudio = true
			r.Languages = append(r.Languages, "en")
			flags["multi"] = true
		}
	}
}
