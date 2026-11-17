package database

import (
	"context"
	"testing"
	"testing/fstest"
)

// A fork renumbers its own migration after the upstream one that took its
// number: the fork's migration is not run again, its record moves to the
// new number, and the upstream migration that now sits below the newest
// applied number still runs.
func TestForkMigrationRenumberedBehindUpstream(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	base := upTo(t, 122)
	withFork := fstest.MapFS{}
	for n, f := range base.(fstest.MapFS) {
		withFork[n] = f
	}
	// Running this twice would fail: the table exists.
	withFork["0123_fork_feature.sql"] = &fstest.MapFile{Data: []byte(`CREATE TABLE fork_feature (id INTEGER PRIMARY KEY);`)}
	s, err := Open(dir, Options{MainMigrations: withFork})
	if err != nil {
		t.Fatal(err)
	}
	_ = s.Close()

	updated := fstest.MapFS{}
	for n, f := range base.(fstest.MapFS) {
		updated[n] = f
	}
	updated["0123_upstream_feature.sql"] = &fstest.MapFile{Data: []byte(`CREATE TABLE upstream_feature (id INTEGER PRIMARY KEY);`)}
	updated["0124_fork_feature.sql"] = withFork["0123_fork_feature.sql"]
	s, err = Open(dir, Options{MainMigrations: updated})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	defer s.Close()
	var n int
	if err := s.Read().QueryRow(`SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name IN ('fork_feature', 'upstream_feature')`).Scan(&n); err != nil || n != 2 {
		t.Fatalf("tables %d %v", n, err)
	}
	rows := map[int]string{}
	r, err := s.Read().QueryContext(ctx, `SELECT version, name FROM schema_migrations WHERE version >= 123`)
	if err != nil {
		t.Fatal(err)
	}
	for r.Next() {
		var v int
		var name string
		_ = r.Scan(&v, &name)
		rows[v] = name
	}
	_ = r.Close()
	if len(rows) != 2 || rows[123] != "0123_upstream_feature.sql" || rows[124] != "0124_fork_feature.sql" {
		t.Fatalf("records %v", rows)
	}
}

// A database with no gaps runs nothing new on a plain restart.
func TestNoMigrationRunsTwice(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < 2; i++ {
		s, err := Open(dir, Options{})
		if err != nil {
			t.Fatalf("open %d: %v", i, err)
		}
		_ = s.Close()
	}
}
