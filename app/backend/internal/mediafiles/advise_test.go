package mediafiles

import (
	"slices"
	"testing"
)

func TestAdviseRules(t *testing.T) {
	german := func(order string, good int) PickRulesBody {
		b := GoodEveryday()
		b.Traits = append(b.Traits, TraitPoints{Trait: "own:1", Points: 11000})
		b.Order, b.GoodEnoughPoints = order, &good
		return b
	}
	if got := AdviseRules(german("quality_then_points", 0)); !slices.Equal(got, []string{AdviceOrder, AdviceGoodEnough}) {
		t.Fatalf("imported German set: %q", got)
	}
	if got := AdviseRules(german("points_then_quality", 11000)); got != nil {
		t.Fatalf("set that ranks German first: %q", got)
	}
	if got := AdviseRules(german("points_then_quality", 35000)); got != nil {
		t.Fatalf("TRaSH's high good enough is its own choice: %q", got)
	}
	empty := german("points_then_quality", 0)
	empty.GoodEnoughPoints = nil
	if got := AdviseRules(empty); !slices.Equal(got, []string{AdviceGoodEnough}) {
		t.Fatalf("an empty good enough counts by the step alone: %q", got)
	}
	plain := GoodEveryday()
	zero := 0
	plain.GoodEnoughPoints = &zero
	if got := AdviseRules(plain); got != nil {
		t.Fatalf("rules without a language trait: %q", got)
	}
	one := german("quality_then_points", 11000)
	for i := range one.Steps {
		one.Steps[i].Allowed = one.Steps[i].Quality == "1080p_web"
	}
	if got := AdviseRules(one); got != nil {
		t.Fatalf("one step has nothing to rank: %q", got)
	}
}
