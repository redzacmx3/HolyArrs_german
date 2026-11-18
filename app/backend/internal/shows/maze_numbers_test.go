package shows

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	domain "github.com/HolyArrs/holyarrs/app/backend/internal/core/model"
	"github.com/HolyArrs/holyarrs/app/backend/internal/mediafiles"
	"github.com/HolyArrs/holyarrs/app/backend/internal/metadata"
)

// Release names by TVmaze: a refresh stores TVmaze's numbers for each
// episode it finds there; while TVmaze cannot be read they stay.
func TestRefreshStoresMazeNumbers(t *testing.T) {
	svc, s, tmdb, maze := addShowService(t)
	ctx := context.Background()
	addShowFolder(t, s)
	tmdb.set(labShow("tmdb", "100"), "tmdb:100", "tvdb:200", "imdb:tt0100")
	m := labShow("tvmaze", "300")
	m.Episodes[3].Season, m.Episodes[3].Number = 1, 3 // "Return" counts on in season 1 on TVmaze
	maze.set(m, "tvmaze:300", "tvdb:200", "imdb:tt0100")
	id, _, err := svc.AddShow(ctx, "tmdb:100", ShowAddOptions{}, domain.System)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Write(ctx, func(tx *sql.Tx) error {
		_, err := tx.Exec("UPDATE titles SET release_numbers = 'tvmaze' WHERE id = ?", id)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	named := func() map[string][2]int {
		out := map[string][2]int{}
		rows, err := s.Read().QueryContext(ctx, `SELECT name, COALESCE(json_extract(numbers, '$.name_season'), 0), COALESCE(json_extract(numbers, '$.name_episode'), 0)
			FROM items WHERE title_id = ? AND kind = 'episode'`, id)
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var n string
			var v [2]int
			if err := rows.Scan(&n, &v[0], &v[1]); err != nil {
				t.Fatal(err)
			}
			out[n] = v
		}
		return out
	}
	if _, err := svc.RefreshShow(ctx, id, domain.System, nil); err != nil {
		t.Fatal(err)
	}
	got := named()
	if got["Return"] != [2]int{1, 3} || got["Pilot"] != [2]int{1, 1} || got["Behind the Lab"] != [2]int{} {
		t.Fatalf("TVmaze numbers %v", got)
	}
	var ti mediafiles.Title
	_ = s.Read().QueryRowContext(ctx, "SELECT release_numbers FROM titles WHERE id = ?", id).Scan(&ti.ReleaseNumbers)
	if ti.ReleaseNames() != mediafiles.ReleaseNumbersTVmaze {
		t.Fatalf("release names %q", ti.ReleaseNames())
	}

	// TVmaze unreachable: the numbers stay.
	maze.err = errors.New("TVmaze is down")
	if _, err := svc.RefreshShow(ctx, id, domain.System, nil); err != nil {
		t.Fatal(err)
	}
	if got := named(); got["Return"] != [2]int{1, 3} {
		t.Fatalf("kept numbers %v", got)
	}
}

// serverMaze is TVmaze as the metadata server serves it: no lookups by TVDB
// or IMDb id, and search hits without those ids.
type serverMaze struct{ *fakeShows }

func (m serverMaze) Show(ctx context.Context, id string) (metadata.Show, error) {
	if !strings.HasPrefix(id, "tvmaze:") {
		return metadata.Show{}, metadata.ErrShowNotFound
	}
	return m.fakeShows.Show(ctx, id)
}

func (m serverMaze) SearchShows(ctx context.Context, q string) ([]metadata.Show, error) {
	hits, err := m.fakeShows.SearchShows(ctx, q)
	for i := range hits {
		hits[i].IDs = map[string]string{"tvmaze": hits[i].IDs["tvmaze"]}
	}
	return hits, err
}

// Through the metadata server the show is found by name and checked by
// reading the hit in full; a show TVmaze does not have gives the reason.
func TestRefreshMazeNumbersThroughServer(t *testing.T) {
	svc, s, tmdb, maze := addShowService(t)
	ctx := context.Background()
	addShowFolder(t, s)
	tmdb.set(labShow("tmdb", "100"), "tmdb:100", "tvdb:200", "imdb:tt0100")
	maze.shows = nil
	maze.set(labShow("tvmaze", "300"), "tvmaze:300")
	svc.Shows["tvmaze"] = serverMaze{maze}
	id, _, err := svc.AddShow(ctx, "tmdb:100", ShowAddOptions{}, domain.System)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Write(ctx, func(tx *sql.Tx) error {
		_, err := tx.Exec("UPDATE titles SET release_numbers = 'tvmaze', source_ids = json_remove(source_ids, '$.tvmaze') WHERE id = ?", id)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	out, err := svc.RefreshShow(ctx, id, domain.System, nil)
	if err != nil {
		t.Fatal(err)
	}
	if out.TVmazeError != "" || out.TVmazeNumbered != 5 {
		t.Fatalf("refresh %+v", out)
	}

	other := labShow("tvmaze", "301")
	other.IDs["tvdb"], other.IDs["imdb"] = "999", "tt0999"
	maze.shows = nil
	maze.set(other, "tvmaze:301")
	out, err = svc.RefreshShow(ctx, id, domain.System, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.TVmazeError, "no show with TVDB id") {
		t.Fatalf("want a reason, got %+v", out)
	}
}
