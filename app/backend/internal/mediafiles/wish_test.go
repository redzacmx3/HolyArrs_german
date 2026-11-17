package mediafiles

import (
	"slices"
	"testing"
)

func allowedSteps(b PickRulesBody) []string {
	var out []string
	for _, s := range b.Steps {
		if s.Allowed {
			out = append(out, s.Quality)
		}
	}
	return out
}

func TestApplyWish(t *testing.T) {
	yes, no := true, false
	// A 1080p set aimed at 720p keeps its 720p steps and is good enough there.
	b := ApplyWish(GoodEveryday(), Wish{Quality: "720p"})
	if !slices.Equal(allowedSteps(b), []string{"720p_tv", "720p_web", "720p_bluray"}) || b.GoodEnoughAt != "720p_bluray" {
		t.Fatalf("720p: %v %s", allowedSteps(b), b.GoodEnoughAt)
	}
	// A 1080p set aimed at 4K gains the 4K web and Blu-ray steps and their sizes.
	b = ApplyWish(GoodEveryday(), Wish{Quality: "2160p"})
	if !slices.Contains(allowedSteps(b), "2160p_web") || slices.Contains(allowedSteps(b), "2160p_remux") || b.GoodEnoughAt != "2160p_bluray" {
		t.Fatalf("2160p: %v %s", allowedSteps(b), b.GoodEnoughAt)
	}
	if b.SizeLimits["2160p_web"].Max <= GoodEveryday().SizeLimits["*"].Max {
		t.Fatalf("4K sizes %+v", b.SizeLimits)
	}
	// A 4K set aimed at 1080p keeps only what it allowed up to 1080p.
	b = ApplyWish(Best(), Wish{Quality: "1080p"})
	if !slices.Equal(allowedSteps(b), []string{"1080p_bluray"}) || b.GoodEnoughAt != "1080p_bluray" {
		t.Fatalf("1080p: %v %s", allowedSteps(b), b.GoodEnoughAt)
	}
	// Languages replace; an inner wish wins over an outer one.
	w := Wish{Languages: []string{"de"}, OriginalOK: &yes, Quality: "1080p"}.Over(Wish{Languages: []string{"en"}, OriginalOK: &no})
	b = ApplyWish(GoodEveryday(), w)
	if !slices.Equal(b.Language.Wanted, []string{"en"}) || b.Language.OriginalOK || b.GoodEnoughAt != "1080p_web" {
		t.Fatalf("languages %+v %s", b.Language, b.GoodEnoughAt)
	}
	if !(Wish{}).Empty() || (Wish{Languages: []string{}}).Empty() {
		t.Fatal("empty")
	}
	if (Wish{Quality: "8k"}).Check() == nil || (Wish{Languages: []string{"DE"}}).Check() == nil || (Wish{Languages: []string{"de", "pt"}}).Check() != nil {
		t.Fatal("check")
	}
}
