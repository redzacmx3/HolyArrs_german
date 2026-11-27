package database

import (
	"context"
	"database/sql"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
)

// Fork: every migration but own_trait_per_minute, found by name, as an
// upstream update can give it a new number. Since the 1.3.6 sync the fork's
// own_trait_language runs after it on a new database; a database from
// before already has it, so it stays in.
func beforePerMinute(t *testing.T) fs.FS {
	t.Helper()
	src, _ := fs.Sub(embedded, "migrations/main")
	ms, err := loadMigrations(src)
	if err != nil {
		t.Fatal(err)
	}
	stop := 0
	for _, m := range ms {
		if strings.HasSuffix(m.name, "_own_trait_per_minute.sql") {
			stop = m.version
		}
	}
	if stop == 0 {
		t.Fatal("own_trait_per_minute migration not found")
	}
	out := fstest.MapFS{}
	for _, m := range ms {
		if m.version != stop {
			out[m.name] = &fstest.MapFile{Data: []byte(m.sql)}
		}
	}
	return out
}

// The fork's own_trait_per_minute rebuilds own_traits to allow the
// size_per_minute kind and keeps every trait, its id, and the AUTOINCREMENT
// mark, so a removed trait's id is never reused.
func TestOwnTraitsKeptByPerMinuteMigration(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	old, err := Open(dir, Options{MainMigrations: beforePerMinute(t)})
	if err != nil {
		t.Fatal(err)
	}
	err = old.Write(ctx, func(tx *sql.Tx) error {
		_, err := tx.Exec(`INSERT INTO own_traits (id, created_at, updated_at, name, match, terms, source, source_id, specs) VALUES
			(4, 'c', 'u', 'Huge', 'size', '["60-"]', '', '', '[]'),
			(5, 'c', 'u', 'German', 'language', '["de"]', '', '', '[]'),
			(8, 'c', 'u', 'German Microsized', 'name', '[]', 'trash', 'sonarr:x', '[{"kind":"group","value":"^(FuN)$"}]')`)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`INSERT INTO own_traits (id, created_at, updated_at, name, match) VALUES (20, 'c', 'u', 'Gone', 'name')`)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`DELETE FROM own_traits WHERE id = 20`)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = old.Close()

	s, err := Open(dir, Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	var n int
	if err := s.Read().QueryRow(`SELECT count(*) FROM own_traits WHERE
		(id = 4 AND match = 'size' AND terms = '["60-"]') OR
		(id = 5 AND match = 'language' AND terms = '["de"]') OR
		(id = 8 AND source = 'trash' AND source_id = 'sonarr:x')`).Scan(&n); err != nil || n != 3 {
		t.Fatalf("traits kept %d %v", n, err)
	}
	var id int64
	err = s.Write(ctx, func(tx *sql.Tx) error {
		r, err := tx.Exec(`INSERT INTO own_traits (created_at, updated_at, name, match, terms) VALUES ('c', 'u', 'Small a minute', 'size_per_minute', '["-10"]')`)
		if err != nil {
			return err
		}
		id, err = r.LastInsertId()
		return err
	})
	if err != nil {
		t.Fatalf("size per minute trait: %v", err)
	}
	if id <= 20 {
		t.Fatalf("a removed trait's id came back: %d", id)
	}
	err = s.Write(ctx, func(tx *sql.Tx) error {
		_, err := tx.Exec(`INSERT INTO own_traits (created_at, updated_at, name, match) VALUES ('c', 'u', 'x', 'title')`)
		return err
	})
	if err == nil {
		t.Fatal("an unknown match was accepted")
	}
}
