package sources

import (
	"context"
	"database/sql"
	"net/url"
	"testing"

	"github.com/HolyArrs/holyarrs/app/backend/internal/core/jobs"
	domain "github.com/HolyArrs/holyarrs/app/backend/internal/core/model"
	"github.com/HolyArrs/holyarrs/app/backend/internal/mediafiles"
)

// A show on a chosen episode order whose release names count by TMDB's own
// numbering (German releases of some kids' shows): with release_numbers
// tmdb the search asks for and takes S01E02 for the episode TMDB numbers
// S01E02, although the order calls it S03E01 and gives S01E02 to another.
func TestSearch_ReleaseNumbersTMDB(t *testing.T) {
	const name = "Copper.Valley.S01E02.1080p.WEB-DL.H.264-QUILL"
	for _, tmdb := range []bool{true, false} {
		ts := &tvServer{answer: func(url.Values) []string { return []string{name} }}
		f := seedShow(t, ts.start(t), nil)
		setAlt(t, f, "S01E02", 3, 1)
		setAlt(t, f, "S01E03", 1, 2)
		numbers := ""
		if tmdb {
			numbers = mediafiles.ReleaseNumbersTMDB
		}
		if err := f.s.Write(context.Background(), func(tx *sql.Tx) error {
			_, err := tx.Exec("UPDATE titles SET episode_group_id = 'grp', release_numbers = ? WHERE id = ?", numbers, f.titleID)
			return err
		}); err != nil {
			t.Fatal(err)
		}
		runKind(t, f.fixture, &Searcher{Store: f.s}, jobs.KindSearchSlot, &domain.Subject{Type: "slot", ID: f.slots["S01E02"]}, nil)
		wantSeason, wantEp, want := "3", "1", mediafiles.SlotWanted
		if tmdb {
			wantSeason, wantEp, want = "1", "2", mediafiles.SlotChosen
		}
		if q := ts.queries[0]; q.Get("season") != wantSeason || q.Get("ep") != wantEp {
			t.Fatalf("tmdb %v: first query %v", tmdb, q)
		}
		if st := f.state(t, "S01E02"); st != want {
			r := relByName(t, f.fixture, name)
			t.Fatalf("tmdb %v: S01E02 %s, want %s (%s %v)", tmdb, st, want, r.State, reasonCodes(r))
		}
	}
}
