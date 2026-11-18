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

// Release names by TVmaze (TVDB-style numbers, as German FuN releases of
// PAW Patrol count): the search asks for the TVmaze numbers and a release
// named so fills the episode, whatever the show's own numbers are.
func TestSearch_ReleaseNumbersTVmaze(t *testing.T) {
	const name = "Copper.Valley.S02E05.1080p.WEB-DL.H.264-FuN"
	ts := &tvServer{answer: func(url.Values) []string { return []string{name} }}
	f := seedShow(t, ts.start(t), nil)
	setAlt(t, f, "S01E02", 3, 1)
	for label, n := range map[string][2]int{"S01E02": {2, 5}, "S01E03": {2, 6}} {
		if err := f.s.Write(context.Background(), func(tx *sql.Tx) error {
			_, err := tx.Exec(`UPDATE items SET numbers = json_set(numbers, '$.name_season', ?, '$.name_episode', ?) WHERE id = ?`, n[0], n[1], f.items[label])
			return err
		}); err != nil {
			t.Fatal(err)
		}
	}
	if err := f.s.Write(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec("UPDATE titles SET episode_group_id = 'grp', release_numbers = ? WHERE id = ?", mediafiles.ReleaseNumbersTVmaze, f.titleID)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	runKind(t, f.fixture, &Searcher{Store: f.s}, jobs.KindSearchSlot, &domain.Subject{Type: "slot", ID: f.slots["S01E02"]}, nil)
	if q := ts.queries[0]; q.Get("season") != "2" || q.Get("ep") != "5" {
		t.Fatalf("first query %v", q)
	}
	if st := f.state(t, "S01E02"); st != mediafiles.SlotChosen {
		r := relByName(t, f.fixture, name)
		t.Fatalf("S01E02 %s (%s %v)", st, r.State, reasonCodes(r))
	}
}

// With release names by TVmaze, an episode search no Source answers asks
// the release-name season, where the double episode is found.
func TestSearch_ReleaseNumbersTVmazeSeasonFallback(t *testing.T) {
	const name = "Copper.Valley.S02E04-E05.1080p.WEB-DL.H.264-FuN"
	ts := &tvServer{answer: func(q url.Values) []string {
		if q.Get("ep") != "" {
			return nil
		}
		return []string{name}
	}}
	f := seedShow(t, ts.start(t), nil)
	for label, n := range map[string][2]int{"S01E02": {2, 5}, "S01E03": {2, 4}} {
		if err := f.s.Write(context.Background(), func(tx *sql.Tx) error {
			_, err := tx.Exec(`UPDATE items SET numbers = json_set(numbers, '$.name_season', ?, '$.name_episode', ?) WHERE id = ?`, n[0], n[1], f.items[label])
			return err
		}); err != nil {
			t.Fatal(err)
		}
	}
	if err := f.s.Write(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec("UPDATE titles SET release_numbers = ? WHERE id = ?", mediafiles.ReleaseNumbersTVmaze, f.titleID)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	runKind(t, f.fixture, &Searcher{Store: f.s}, jobs.KindSearchSlot, &domain.Subject{Type: "slot", ID: f.slots["S01E02"]}, nil)
	last := ts.queries[len(ts.queries)-1]
	if len(ts.queries) < 2 || last.Get("season") != "2" || last.Get("ep") != "" {
		t.Fatalf("queries %v", ts.queries)
	}
	if st := f.state(t, "S01E02"); st != mediafiles.SlotChosen {
		r := relByName(t, f.fixture, name)
		t.Fatalf("S01E02 %s (%s %v)", st, r.State, reasonCodes(r))
	}
}
