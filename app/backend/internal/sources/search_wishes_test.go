package sources

import (
	"testing"

	"github.com/HolyArrs/holyarrs/app/backend/internal/core/database"
	"github.com/HolyArrs/holyarrs/app/backend/internal/core/releasename"
	"github.com/HolyArrs/holyarrs/app/backend/internal/mediafiles"
)

// TestShowWishes: the show's wish, the season's and the episode's lay over
// each other, and a pack that holds episodes with different wishes must
// pass each of them.
func TestShowWishes(t *testing.T) {
	part := int64(7)
	wishes := []mediafiles.PickWish{
		{Wish: mediafiles.Wish{Languages: []string{"de"}}}, // every show
		{TitleID: 1, Wish: mediafiles.Wish{Quality: "1080p"}},
		{TitleID: 1, PartID: part, Wish: mediafiles.Wish{Quality: "720p"}},
		{TitleID: 1, ItemID: 3, Wish: mediafiles.Wish{Languages: []string{"en"}}},
	}
	pr := mediafiles.PickRules{Body: mediafiles.GoodEveryday()}
	_, e1 := wishOf(wishes, mediafiles.Item{Meta: database.Meta{ID: 1}}, pr)
	_, e2 := wishOf(wishes, mediafiles.Item{Meta: database.Meta{ID: 2}, PartID: &part}, pr)
	_, e3 := wishOf(wishes, mediafiles.Item{Meta: database.Meta{ID: 3}, PartID: &part}, pr)
	if e1.Body.GoodEnoughAt != "1080p_web" || e1.Body.Language.Wanted[0] != "de" {
		t.Fatalf("show level %s %v", e1.Body.GoodEnoughAt, e1.Body.Language)
	}
	if e2.Body.GoodEnoughAt != "720p_bluray" || e2.Body.Language.Wanted[0] != "de" {
		t.Fatalf("season level %s %v", e2.Body.GoodEnoughAt, e2.Body.Language)
	}
	if e3.Body.GoodEnoughAt != "720p_bluray" || e3.Body.Language.Wanted[0] != "en" {
		t.Fatalf("episode level %s %v", e3.Body.GoodEnoughAt, e3.Body.Language)
	}
	if k, _ := wishOf(nil, mediafiles.Item{}, pr); k != "" {
		t.Fatalf("no wish gave key %q", k)
	}

	k2, _ := wishOf(wishes, mediafiles.Item{Meta: database.Meta{ID: 2}, PartID: &part}, pr)
	k3, _ := wishOf(wishes, mediafiles.Item{Meta: database.Meta{ID: 3}, PartID: &part}, pr)
	w := showWishes{base: pr.Body, bySlot: map[int64]string{20: k2, 30: k3},
		bodies: map[string]mediafiles.PickRulesBody{k2: e2.Body, k3: e3.Body}}
	title := &mediafiles.TitleRef{ID: 1, Name: "Show", Kind: "show", OriginalLanguage: "en"}
	name := "Show.S01.German.720p.WEB.H264-GRP"
	rel := mediafiles.ReleaseFacts{Name: name, Parse: releasename.Parse(name, releasename.Hints{Show: true}), Size: 2 << 30, Protocol: "usenet"}
	ss := mediafiles.SlotSet{Match: mediafiles.Match{Title: title, Conf: releasename.Sure}, Out: true, RuntimeMinutes: 90, SlotIDs: []int64{20}}
	if v := w.judge(rel, ss, pr.Body); !v.Passed {
		t.Fatalf("German pack for the German episode refused: %+v", v.Reasons)
	}
	ss.SlotIDs = []int64{20, 30}
	if v := w.judge(rel, ss, pr.Body); v.Passed {
		t.Fatal("German pack passed although episode 3 wants English only")
	}
}
