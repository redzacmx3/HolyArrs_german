package library

import (
	"context"
	"database/sql"
	"slices"
	"testing"
	"time"

	domain "github.com/HolyArrs/holyarrs/app/backend/internal/core/model"
)

func TestUpkeepNeverWatched(t *testing.T) {
	svc, s, lf, id := renameWorld(t)
	ctx := context.Background()
	f := putFile(t, s, lf, id, "Dune - Part Two (2024)/Dune - Part Two (2024) 1080p BluRay.mkv", "Dune.Part.Two.2024.1080p.BluRay-HOLYARRS", "imported")
	later := time.Now().AddDate(0, 0, 400)
	svc.Now = func() time.Time { return later }
	watched, history := map[int64]bool{}, false
	svc.WatchedTitles = func(context.Context) (map[int64]bool, bool, error) { return watched, history, nil }

	l, err := svc.Upkeep(ctx)
	if err != nil || len(l.Items) != 0 || l.WatchHistory {
		t.Fatalf("no watch history must name nothing: %+v %v", l, err)
	}
	history, watched[id+1000] = true, true
	l, _ = svc.Upkeep(ctx)
	if len(l.Items) != 1 || l.Items[0].TitleID != id || l.Items[0].Rule != RuleNeverWatched || l.Items[0].Mode != UpkeepSuggest || l.Items[0].Files != 1 {
		t.Fatalf("unwatched movie: %+v", l.Items)
	}
	svc.Now = func() time.Time { return time.Now().AddDate(0, 0, 100) }
	if l, _ = svc.Upkeep(ctx); len(l.Items) != 0 {
		t.Fatal("a movie younger than the days is kept")
	}
	svc.Now = func() time.Time { return later }
	watched[id] = true
	if l, _ = svc.Upkeep(ctx); len(l.Items) != 0 {
		t.Fatal("a watched movie is kept")
	}
	delete(watched, id)

	// The daily run leaves suggestions alone.
	if done, err := svc.ApplyUpkeep(ctx, nil, true, domain.System); err != nil || len(done.Removed) != 0 {
		t.Fatalf("suggest mode removed: %+v %v", done, err)
	}
	u := DefaultUpkeep()
	u.NeverWatched = "sometimes"
	if err := svc.SaveUpkeep(ctx, u, domain.System); err == nil {
		t.Fatal("a bad mode was saved")
	}
	u.NeverWatched, u.Paused = UpkeepAuto, true
	if err := svc.SaveUpkeep(ctx, u, domain.System); err != nil {
		t.Fatal(err)
	}
	if done, _ := svc.ApplyUpkeep(ctx, nil, true, domain.System); len(done.Removed) != 0 {
		t.Fatal("paused upkeep removed")
	}
	u.Paused = false
	_ = svc.SaveUpkeep(ctx, u, domain.System)
	done, err := svc.ApplyUpkeep(ctx, nil, true, domain.System)
	if err != nil || !slices.Equal(done.Removed, []int64{id}) {
		t.Fatalf("auto run: %+v %v", done, err)
	}
	if onDisk(lf, f.RelativePath) || fileRow(t, s, f.ID).State != "gone" {
		t.Fatal("the file did not go to the holding folder")
	}
}

func TestUpkeepNoFileAndOnlyListedTitles(t *testing.T) {
	svc, s, _, id := renameWorld(t)
	ctx := context.Background()
	l, _ := svc.Upkeep(ctx)
	if len(l.Items) != 0 {
		t.Fatalf("a movie kept up to date is still wanted: %+v", l.Items)
	}
	if err := s.Write(ctx, func(tx *sql.Tx) error {
		_, err := tx.Exec("UPDATE versions SET keep_up_to_date = 0 WHERE title_id = ?", id)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	l, _ = svc.Upkeep(ctx)
	if len(l.Items) != 1 || l.Items[0].Rule != RuleNoFile {
		t.Fatalf("no file, nothing keeps it: %+v", l.Items)
	}
	done, err := svc.ApplyUpkeep(ctx, []int64{id + 1000, id}, false, domain.System)
	if err != nil || !slices.Equal(done.Removed, []int64{id}) || !slices.Equal(done.Skipped, []int64{id + 1000}) {
		t.Fatalf("apply: %+v %v", done, err)
	}
}
