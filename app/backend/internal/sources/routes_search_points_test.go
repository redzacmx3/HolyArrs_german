package sources

import (
	"testing"

	domain "github.com/HolyArrs/holyarrs/app/backend/internal/core/model"
)

// Fork: each trait reason carries its points, so "Release auswählen" can list
// what added and what took points away.
func TestReleaseViewTraitPoints(t *testing.T) {
	v := releaseView(domain.Release{Name: "PAW.Patrol.S02E31.German.1080p.WebHD.x265-FuN", Verdict: map[string]any{
		"points": float64(-23300),
		"reasons": []any{
			map[string]any{"step": "points", "rule": "trait:own:7", "effect": "subtract", "detail": "German Microsized: fits", "points": float64(-35000)},
			map[string]any{"step": "points", "rule": "trait:own:2", "effect": "add", "detail": "German: fits", "points": float64(10000)},
			map[string]any{"step": "points", "rule": "below_min", "effect": "fail", "detail": "-23300 points, needs at least 11000"},
		},
	}})
	if len(v.Reasons) != 3 {
		t.Fatalf("reasons = %+v", v.Reasons)
	}
	if r := v.Reasons[0]; r.Points != -35000 || r.Text != "German Microsized: fits (-35000 points)" {
		t.Errorf("subtract reason = %+v", r)
	}
	if r := v.Reasons[1]; r.Points != 10000 {
		t.Errorf("add reason = %+v", r)
	}
	if r := v.Reasons[2]; r.Points != 0 {
		t.Errorf("fail reason has points: %+v", r)
	}
}
