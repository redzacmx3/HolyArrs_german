package releasename

import "slices"

// subWords after a language word say the language is the subtitles'. The
// upstream reading already keeps "sub" and "subs" languages out of the audio
// (Sonarr reads on after them), so only subbed and subtitled take one out.
var subWords = setOf("sub", "subs", "subbed", "subtitled")

// readGermanTags runs after the upstream language reading and adds what the
// German Usenet names need: a language word before a subtitle word names
// subtitles ("German.Subbed", "GERMAN.SUBS"), OmU is the original audio with
// German subtitles, and German DL/ML (read upstream as the "original" flag)
// set OriginalAudio, German ML also multi.
func readGermanTags(r *Result, low []string, flags map[string]bool, has func(...string) bool) {
	for i, t := range low {
		if i+1 >= len(low) || !subWords[low[i+1]] {
			continue
		}
		code, ok := langWords[t]
		if !ok {
			code, ok = moreLangWords[t]
		}
		if !ok {
			continue
		}
		if w := low[i+1]; w == "subbed" || w == "subtitled" {
			r.Languages = slices.DeleteFunc(r.Languages, func(l string) bool { return l == code })
		}
		if !slices.Contains(r.SubLanguages, code) {
			r.SubLanguages = append(r.SubLanguages, code)
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
	if len(r.Languages) > 0 && r.Languages[0] == "de" && flags["original"] {
		r.OriginalAudio = true
		if has("ml") {
			flags["multi"] = true
		}
	}
}
