package plugins

import (
	"context"
	"testing"
)

// The upstream tests install German as a pack, so they run with no fork
// built-in languages; TestForkGermanIsBuiltIn turns German back on.
func init() { builtInLanguages = nil }

// German ships in the fork's build: it is listed as built in, offered on
// the sign-in page and never fetched as a pack.
func TestForkGermanIsBuiltIn(t *testing.T) {
	builtInLanguages = []string{"de"}
	t.Cleanup(func() { builtInLanguages = nil })
	l := &Languages{Dir: t.TempDir()}
	rows := l.List()
	if len(rows) < 2 || rows[1].Code != "de" || rows[1].State != "built_in" {
		t.Fatalf("rows %+v", rows[:2])
	}
	for _, r := range rows[2:] {
		if r.Code == "de" {
			t.Fatal("German listed twice")
		}
	}
	if err := l.Install(context.Background(), "de"); err != nil {
		t.Fatalf("install de: %v", err)
	}
	if knownLanguage("de") {
		t.Fatal("German would be fetched as a pack")
	}
}
