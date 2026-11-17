package migrate

import "testing"

func TestOpenQualitiesSeesMergedGroups(t *testing.T) {
	quality := map[int64]string{3: "WEBDL-1080p", 7: "Bluray-1080p", 5: "WEBDL-720p"}
	group := func(id int64, allowed bool, qs ...int64) map[string]any {
		var sub []any
		for _, q := range qs {
			sub = append(sub, map[string]any{"quality": q, "allowed": allowed})
		}
		return map[string]any{"id": id, "name": "Merged QPs", "allowed": allowed, "items": sub}
	}
	single := func(q int64, allowed bool) map[string]any {
		return map[string]any{"quality": q, "allowed": allowed, "items": []any{}}
	}
	allowed, cut, merged := openQualities([]any{group(1001, true, 5, 3, 7)}, quality, 1001)
	if !merged || len(allowed) != 3 || len(cut) != 3 {
		t.Fatalf("merged group: allowed %v cut %v merged %v", allowed, cut, merged)
	}
	_, _, merged = openQualities([]any{single(3, true), single(7, true), group(1001, false, 5, 3)}, quality, 7)
	if merged {
		t.Fatal("separate qualities and a refused group are not merged")
	}
}
