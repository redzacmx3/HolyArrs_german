package health

import (
	"strings"
	"testing"

	"github.com/HolyArrs/holyarrs/app/backend/internal/mediafiles"
)

// A used set whose German trait loses on quality gets a warning naming it;
// an unused one stays quiet, and the warning clears once points lead.
func TestPickRulesCheck(t *testing.T) {
	w := newWorld(t)
	body := `{"steps":[{"quality":"1080p_web","allowed":true},{"quality":"1080p_bluray","allowed":true}],"good_enough_at":"1080p_bluray",
		"traits":[{"trait":"own:1","points":11000}],"must_have":[],"never":[],"min_points":0,"good_enough_points":0,"size_limits":{},
		"protocol":{},"language":{"wanted":["de"],"original_ok":true},"order":"quality_then_points","upgrade_min_points":1}`
	for _, n := range []string{"HD German", "Unused German"} {
		w.exec(t, `INSERT INTO pick_rules (created_at, created_by, updated_at, updated_by, name, kinds, body)
			VALUES ('2026-10-01T00:00:00Z', '{}', '2026-10-01T00:00:00Z', '{}', ?, '["movie"]', ?)`, n, body)
	}
	w.exec(t, `INSERT INTO titles (created_at, created_by, updated_at, updated_by, kind, name, sort_name, year, source_ids, metadata_source, why_added, state)
		VALUES ('2026-10-01T00:00:00Z', '{}', '2026-10-01T00:00:00Z', '{}', 'movie', 'Film', 'Film', 2020, '{}', 'tmdb', 'by_hand', 'active')`)
	w.exec(t, `INSERT INTO versions (created_at, created_by, updated_at, updated_by, title_id, label, pick_rules_id, library_folder_id, folder_name, is_main)
		VALUES ('2026-10-01T00:00:00Z', '{}', '2026-10-01T00:00:00Z', '{}', (SELECT id FROM titles WHERE name = 'Film'), 'Main',
		(SELECT id FROM pick_rules WHERE name = 'HD German'), 1, 'Film (2020)', 1)`)
	var used, unused string
	for k, c := range w.run(t) {
		if c.Kind != "pick_rules" {
			continue
		}
		switch c.Subject {
		case "HD German":
			used = k
			if c.State != Warning || !strings.Contains(c.Message, mediafiles.AdviceOrder) || !strings.Contains(c.Message, mediafiles.AdviceGoodEnough) {
				t.Fatalf("used German set: %+v", c)
			}
		case "Unused German":
			unused = k
		}
	}
	if used == "" || unused != "" {
		t.Fatalf("used %q, unused %q: only used sets are checked", used, unused)
	}
	w.exec(t, `UPDATE pick_rules SET body = json_set(json_set(body, '$.order', 'points_then_quality'), '$.good_enough_points', 11000) WHERE name = 'HD German'`)
	if c := w.run(t)[used]; c.State != OK || c.Message != "" {
		t.Fatalf("after the fix: %+v", c)
	}
}
