package mediafiles

import (
	"strings"
	"testing"
	"time"

	"github.com/HolyArrs/holyarrs/app/backend/internal/core/database"
	"github.com/HolyArrs/holyarrs/app/backend/internal/core/releasename"
)

// TestSizePerMinuteTrait (fork): the size is read against the runtime the
// Release covers, so a short kids' episode is not called small for its size
// alone; no size or no runtime never fits.
func TestSizePerMinuteTrait(t *testing.T) {
	t.Cleanup(func() { SetOwnTraits(nil) })
	SetOwnTraits([]OwnTrait{
		{Meta: database.Meta{ID: 1}, Name: "Small a minute", Match: MatchSizePerMinute, Terms: []string{"-10"}},
		{Meta: database.Meta{ID: 2}, Name: "Middle a minute", Match: MatchSizePerMinute, Terms: []string{"10-40"}},
		{Meta: database.Meta{ID: 3}, Name: "Broken", Match: MatchSizePerMinute, Terms: []string{"x"}},
	})
	const mb = int64(1 << 20)
	cases := []struct {
		id      int64
		size    int64
		runtime int
		want    bool
	}{
		{1, 377 * mb, 12, false}, // PAW Patrol FuN: 31 MB a minute
		{1, 100 * mb, 12, true},  // 8 MB a minute
		{1, 377 * mb, 0, false},  // runtime not known
		{1, 0, 12, false},        // size not known
		{2, 377 * mb, 12, true},
		{2, 1400 * mb, 45, true},
		{2, 2000 * mb, 45, false},
		{3, 100 * mb, 12, false},
	}
	for _, c := range cases {
		p := releasename.Parse("PAW.Patrol.S02E31.German.1080p.WebHD.x265-FuN", releasename.Hints{})
		p.Size = c.size
		if got, why := ownMatchesAt(c.id, p, c.runtime); got != c.want {
			t.Errorf("trait %d at %d MB for %d min = %v (%s), want %v", c.id, c.size/mb, c.runtime, got, why, c.want)
		}
	}
	p := releasename.Parse("PAW.Patrol.S02E31.German.1080p.WebHD.x265-FuN", releasename.Hints{})
	p.Size = 377 * mb
	if _, why := ownMatchesAt(2, p, 12); !strings.Contains(why, "31.4 MB a minute (377 MB for 12 min)") {
		t.Errorf("why = %q", why)
	}
	// Points use the runtime too.
	rules := PickRulesBody{Traits: []TraitPoints{{Trait: OwnTraitKey(1), Points: -35000}, {Trait: OwnTraitKey(2), Points: 100}}}
	if n, _ := scoreTraitsAt(p, rules, 12); n != 100 {
		t.Errorf("points at 12 min = %d, want 100", n)
	}
	if n, _ := scoreTraitsAt(p, rules, 60); n != -35000 {
		t.Errorf("points at 60 min = %d, want -35000", n)
	}
}

func TestCheckSizePerMinuteTrait(t *testing.T) {
	for _, v := range []string{"-10", " 10-40 ", "2.5-"} {
		tr := OwnTrait{Name: "S", Match: MatchSizePerMinute, Terms: []string{v}}
		if err := CheckOwnTrait(&tr); err != nil {
			t.Errorf("%q: %v", v, err)
		}
	}
	for _, v := range []string{"-", "40-10", "x-", "5", "-0"} {
		tr := OwnTrait{Name: "S", Match: MatchSizePerMinute, Terms: []string{v}}
		if err := CheckOwnTrait(&tr); err == nil {
			t.Errorf("%q passed", v)
		}
	}
}

// TestCurrentFileSizePerMinute (fork): the file in place is scored with the
// same runtime as the Release, so a size per minute trait cannot make every
// new Release look better than it.
func TestCurrentFileSizePerMinute(t *testing.T) {
	t.Cleanup(func() { SetOwnTraits(nil) })
	SetOwnTraits([]OwnTrait{{Meta: database.Meta{ID: 1}, Name: "Fine a minute", Match: MatchSizePerMinute, Terms: []string{"20-"}}})
	const mb = int64(1 << 20)
	rules := PickRulesBody{Traits: []TraitPoints{{Trait: OwnTraitKey(1), Points: 500}}}
	cur := CurrentOf("PAW.Patrol.S02E31.German.1080p.WEB.x264-OLD", 377*mb, time.Time{}, rules)
	if cur.Points != 0 {
		t.Fatalf("without runtime the file has %d points", cur.Points)
	}
	rel := releasename.Parse("PAW.Patrol.S02E31.German.1080p.WebHD.x265-FuN", releasename.Hints{})
	if got := currentWithTitle(&cur, rel, rules, 12); got.Points != 500 {
		t.Errorf("file at 12 min = %d points, want 500", got.Points)
	}
	if got := currentWithTitle(&cur, rel, PickRulesBody{}, 12); got != &cur {
		t.Error("rules without the trait scored the file again")
	}
}
