package plugins

import "slices"

// builtInLanguages ship in this fork's build besides English: the fork adds
// its own lines, so the upstream packs never fit its key list and would
// leave the whole app in English. These are never fetched as packs.
var builtInLanguages = []string{"de"}

func builtIn(code string) bool { return slices.Contains(builtInLanguages, code) }
