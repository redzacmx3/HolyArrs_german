package mediafiles

import (
	"errors"
	"slices"
	"strings"
)

// Wish is a small change to a Version's Pick rules for one movie, show,
// season or episode: the languages wanted and the resolution aimed for. An
// empty field keeps what the level above says (show, then season, then
// episode; the Pick rules set is the top).
type Wish struct {
	// Languages replaces the wanted languages (nil: keep). An empty, non-nil
	// list with OriginalOK means the original language only.
	Languages  []string `json:"languages"`
	OriginalOK *bool    `json:"original_ok,omitempty"`
	// Quality is the resolution aimed for: 720p, 1080p or 2160p ("": keep).
	Quality string `json:"quality,omitempty"`
}

// WishQualities are the resolutions a Wish may aim for, low to high.
var WishQualities = []string{"720p", "1080p", "2160p"}

// ErrBadWish says a Wish names a resolution or language that is not known.
var ErrBadWish = errors.New("a wish names an unknown quality or language")

// Empty reports whether the Wish changes nothing.
func (w Wish) Empty() bool {
	return w.Languages == nil && w.OriginalOK == nil && w.Quality == ""
}

// Check makes sure the Wish names known resolutions and language codes.
func (w Wish) Check() error {
	if w.Quality != "" && !slices.Contains(WishQualities, w.Quality) {
		return ErrBadWish
	}
	for _, l := range w.Languages {
		if len(l) < 2 || len(l) > 3 || strings.ToLower(l) != l || strings.Trim(l, "abcdefghijklmnopqrstuvwxyz") != "" {
			return ErrBadWish
		}
	}
	return nil
}

// Over lays inner (a lower level) on w: what inner names wins.
func (w Wish) Over(inner Wish) Wish {
	if inner.Languages != nil {
		w.Languages = inner.Languages
	}
	if inner.OriginalOK != nil {
		w.OriginalOK = inner.OriginalOK
	}
	if inner.Quality != "" {
		w.Quality = inner.Quality
	}
	return w
}

// stepHeight is how high a ladder step's resolution is: 0 below 720p, then
// 1 for 720p, 2 for 1080p and 3 for 2160p.
func stepHeight(step string) int {
	switch {
	case strings.HasPrefix(step, "720p"):
		return 1
	case strings.HasPrefix(step, "1080p"):
		return 2
	case strings.HasPrefix(step, "2160p"):
		return 3
	}
	return 0
}

// ApplyWish returns the rules with the Wish laid on. A language wish
// replaces the wanted languages. A quality wish takes away the steps above
// that resolution, allows its web and Blu-ray steps when the rules allow none
// there, and makes the best of them good enough, so a 4K set can aim for
// 1080p and a 1080p set for 4K.
func ApplyWish(b PickRulesBody, w Wish) PickRulesBody {
	if w.Languages != nil {
		b.Language.Wanted = slices.Clone(w.Languages)
	}
	if w.OriginalOK != nil {
		b.Language.OriginalOK = *w.OriginalOK
	}
	h := slices.Index(WishQualities, w.Quality) + 1
	if h == 0 {
		return b
	}
	steps := slices.Clone(b.Steps)
	has := false
	for i := range steps {
		if stepHeight(steps[i].Quality) > h {
			steps[i].Allowed = false
		}
		has = has || (steps[i].Allowed && stepHeight(steps[i].Quality) == h)
	}
	if !has {
		for i := range steps {
			q := steps[i].Quality
			if stepHeight(q) == h && (strings.HasSuffix(q, "_web") || strings.HasSuffix(q, "_bluray")) {
				steps[i].Allowed = true
			}
		}
		if h == 3 {
			// 4K files are bigger than a 1080p set's limits allow.
			b.SizeLimits = wishSizes(b.SizeLimits, FourK().SizeLimits["*"], h)
		}
	}
	b.Steps = steps
	if stepHeight(b.GoodEnoughAt) != h || !stepAllowed(steps, b.GoodEnoughAt) {
		for _, s := range steps {
			if s.Allowed && stepHeight(s.Quality) == h && !strings.HasSuffix(s.Quality, "_remux") {
				b.GoodEnoughAt = s.Quality
			}
		}
	}
	return b
}

func stepAllowed(steps []Step, q string) bool {
	return slices.ContainsFunc(steps, func(s Step) bool { return s.Quality == q && s.Allowed })
}

// wishSizes gives every step of height h its own size limit l and keeps the
// rest.
func wishSizes(in map[string]SizeLimit, l SizeLimit, h int) map[string]SizeLimit {
	out := make(map[string]SizeLimit, len(in)+3)
	for k, v := range in {
		out[k] = v
	}
	for _, q := range Ladder {
		if stepHeight(q) == h {
			out[q] = l
		}
	}
	return out
}
