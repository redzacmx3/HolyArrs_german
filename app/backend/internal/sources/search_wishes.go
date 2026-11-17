package sources

import (
	"encoding/json"

	"github.com/HolyArrs/holyarrs/app/backend/internal/mediafiles"
)

// wishOf lays the title's, the season's and the episode's pick wishes on the
// rules. key is "" when no wish holds for the episode.
func wishOf(wishes []mediafiles.PickWish, it mediafiles.Item, pr mediafiles.PickRules) (key string, out mediafiles.PickRules) {
	var w mediafiles.Wish
	for _, level := range []func(mediafiles.PickWish) bool{
		func(p mediafiles.PickWish) bool { return p.PartID == 0 && p.ItemID == 0 },
		func(p mediafiles.PickWish) bool { return it.PartID != nil && p.PartID == *it.PartID && p.ItemID == 0 },
		func(p mediafiles.PickWish) bool { return p.ItemID == it.ID },
	} {
		for _, p := range wishes {
			if level(p) {
				w = w.Over(p.Wish)
			}
		}
	}
	if w.Empty() {
		return "", pr
	}
	b, _ := json.Marshal(w)
	pr.Body = mediafiles.ApplyWish(pr.Body, w)
	return string(b), pr
}

// showWishes are the rules each episode of a show search is judged by.
type showWishes struct {
	base   mediafiles.PickRulesBody
	bodies map[string]mediafiles.PickRulesBody
	bySlot map[int64]string
}

func (set showSet) wished() showWishes {
	w := showWishes{base: withVersionAudio(set.pr.Body, set.ver.AudioLanguages), bodies: map[string]mediafiles.PickRulesBody{}, bySlot: map[int64]string{}}
	for _, e := range set.eps {
		if e.sl == nil {
			continue
		}
		w.bySlot[e.sl.ID] = e.wish
		if e.wish != "" {
			w.bodies[e.wish] = withVersionAudio(e.pr.Body, set.ver.AudioLanguages)
		}
	}
	return w
}

// target is the wish every searched episode shares, or "" when they differ
// or have none.
func (w showWishes) target(set showSet) string {
	keys := map[string]bool{}
	for id := range set.target {
		keys[w.bySlot[id]] = true
	}
	if len(keys) != 1 {
		return ""
	}
	for k := range keys {
		return k
	}
	return ""
}

func (w showWishes) body(key string) mediafiles.PickRulesBody {
	if b, ok := w.bodies[key]; ok {
		return b
	}
	return w.base
}

// judge judges a Release by the rules of the episodes it holds. When they
// hold different wishes, it must pass each, and the weakest verdict counts.
func (w showWishes) judge(rel mediafiles.ReleaseFacts, ss mediafiles.SlotSet, fallback mediafiles.PickRulesBody) mediafiles.Verdict {
	keys := []string{}
	seen := map[string]bool{}
	for _, id := range ss.SlotIDs {
		if k := w.bySlot[id]; !seen[k] {
			seen[k] = true
			keys = append(keys, k)
		}
	}
	if len(keys) == 0 {
		return mediafiles.Judge(rel, ss, fallback)
	}
	var worst mediafiles.Verdict
	for i, k := range keys {
		v := mediafiles.Judge(rel, ss, w.body(k))
		if i == 0 || (worst.Passed && !v.Passed) || (worst.Passed == v.Passed && v.Points < worst.Points) {
			worst = v
		}
	}
	return worst
}
