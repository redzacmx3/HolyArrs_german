package database

import (
	"context"
	"strings"
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

// The 1.3.6 sync moved the fork's migrations from 0160-0164 behind
// upstream's new 0160-0162. A database made with the 1.3.0 fork layout runs
// upstream's three and none of the fork's again, and a new database keeps
// the size_per_minute trait kind although own_trait_language now runs after
// own_trait_per_minute.
func TestFork136Renumbering(t *testing.T) {
	all := upTo(t, 1<<30).(fstest.MapFS)
	old := upTo(t, 159).(fstest.MapFS)
	for oldName, now := range map[string]string{
		"0160_own_trait_language.sql":   "0165_own_trait_language.sql",
		"0161_library_upkeep.sql":       "0166_library_upkeep.sql",
		"0162_pick_wishes.sql":          "0167_pick_wishes.sql",
		"0163_release_numbers.sql":      "0163_release_numbers.sql",
		"0164_own_trait_per_minute.sql": "0164_own_trait_per_minute.sql",
	} {
		f, ok := all[now]
		if !ok {
			t.Fatalf("no %s", now)
		}
		old[oldName] = f
	}
	dir := t.TempDir()
	s, err := Open(dir, Options{MainMigrations: old})
	if err != nil {
		t.Fatal(err)
	}
	_ = s.Close()
	s, err = Open(dir, Options{MainMigrations: all})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	defer s.Close()
	var n int
	if err := s.Read().QueryRow(`SELECT count(*) FROM schema_migrations WHERE name IN ('0160_plugin_download_types.sql', '0161_sdr_4k_presets.sql', '0162_lidarr_side_keys.sql', '0165_own_trait_language.sql', '0166_library_upkeep.sql', '0167_pick_wishes.sql')`).Scan(&n); err != nil || n != 6 {
		t.Fatalf("records %d %v", n, err)
	}

	fresh, err := Open(t.TempDir(), Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer fresh.Close()
	var ddl string
	if err := fresh.Read().QueryRow(`SELECT sql FROM sqlite_master WHERE name = 'own_traits'`).Scan(&ddl); err != nil || !strings.Contains(ddl, "'size_per_minute'") || !strings.Contains(ddl, "'language'") {
		t.Fatalf("own_traits on a new database: %v %s", err, ddl)
	}
}
