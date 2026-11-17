package library

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/HolyArrs/holyarrs/app/backend/internal/core/database"
	"github.com/HolyArrs/holyarrs/app/backend/internal/core/jobs"
	domain "github.com/HolyArrs/holyarrs/app/backend/internal/core/model"
	"github.com/HolyArrs/holyarrs/app/backend/internal/mediafiles"
)

// KindUpkeep runs the upkeep rules set to auto (daily).
const KindUpkeep = "library_upkeep"

// Upkeep rule modes: auto acts on its own, suggest only lists, off skips.
const (
	UpkeepAuto    = "auto"
	UpkeepSuggest = "suggest"
	UpkeepOff     = "off"
)

// Upkeep rules.
const (
	RuleNeverWatched = "never_watched" // a movie nobody watched since it came
	RuleNoFile       = "no_file"       // a movie with no file that nothing keeps up to date
)

// The upkeep settings live beside the other where-files-go rules, each
// key starting with upkeep_.
const upkeepGroup, upkeepKey = "where_files_go", "upkeep_"

// UpkeepSettings are the upkeep rules' modes. Paused turns every auto rule
// into a suggestion until it is lifted.
type UpkeepSettings struct {
	NeverWatched     string `json:"never_watched"`
	NeverWatchedDays int    `json:"never_watched_days"`
	NoFile           string `json:"no_file"`
	Paused           bool   `json:"paused"`
}

// DefaultUpkeep is every rule as a suggestion: nothing leaves on its own
// until the owner says so.
func DefaultUpkeep() UpkeepSettings {
	return UpkeepSettings{NeverWatched: UpkeepSuggest, NeverWatchedDays: 365, NoFile: UpkeepSuggest}
}

// ErrUpkeepSettings says a mode or the days are out of range.
var ErrUpkeepSettings = errors.New("never_watched and no_file must be auto, suggest or off, and never_watched_days 30 to 3650")

func (u UpkeepSettings) check() error {
	modes := []string{UpkeepAuto, UpkeepSuggest, UpkeepOff}
	if !slices.Contains(modes, u.NeverWatched) || !slices.Contains(modes, u.NoFile) || u.NeverWatchedDays < 30 || u.NeverWatchedDays > 3650 {
		return ErrUpkeepSettings
	}
	return nil
}

func (u UpkeepSettings) mode(rule string) string {
	m := map[string]string{RuleNeverWatched: u.NeverWatched, RuleNoFile: u.NoFile}[rule]
	if m == UpkeepAuto && u.Paused {
		return UpkeepSuggest
	}
	return m
}

// ReadUpkeep reads the upkeep settings, defaults filling the gaps.
func ReadUpkeep(ctx context.Context, q database.Lister) (UpkeepSettings, error) {
	u := DefaultUpkeep()
	err := database.EachSetting(ctx, q, upkeepGroup, func(k, v string) {
		k, ok := strings.CutPrefix(k, upkeepKey)
		if !ok {
			return
		}
		switch k {
		case RuleNeverWatched:
			u.NeverWatched = v
		case RuleNoFile:
			u.NoFile = v
		case "never_watched_days":
			if n, err := strconv.Atoi(v); err == nil {
				u.NeverWatchedDays = n
			}
		case "paused":
			u.Paused = v == "true"
		}
	})
	if u.check() != nil {
		return DefaultUpkeep(), err
	}
	return u, err
}

// SaveUpkeep stores the upkeep settings.
func (s *Service) SaveUpkeep(ctx context.Context, u UpkeepSettings, by domain.Actor) error {
	if err := u.check(); err != nil {
		return err
	}
	now := s.now()
	return s.Store.Write(ctx, func(tx *sql.Tx) error {
		for k, v := range map[string]string{RuleNeverWatched: u.NeverWatched, RuleNoFile: u.NoFile,
			"never_watched_days": strconv.Itoa(u.NeverWatchedDays), "paused": strconv.FormatBool(u.Paused)} {
			if err := database.PutSetting(tx, upkeepGroup, upkeepKey+k, v, by, now); err != nil {
				return err
			}
		}
		return nil
	})
}

// UpkeepItem is one movie a rule names, with its mode.
type UpkeepItem struct {
	Rule string `json:"rule"`
	Mode string `json:"mode"`
	mediafiles.MovieUpkeep
}

// UpkeepList is what the upkeep rules find now.
type UpkeepList struct {
	Settings UpkeepSettings `json:"settings"`
	// WatchHistory says the media servers' watch history holds anything;
	// without it the never watched rule finds nothing.
	WatchHistory bool         `json:"watch_history"`
	Items        []UpkeepItem `json:"items"`
}

// Upkeep lists the movies each rule that is not off names. A movie named
// by both rules is listed once, under no file.
func (s *Service) Upkeep(ctx context.Context) (UpkeepList, error) {
	db := s.Store.Read()
	out := UpkeepList{Items: []UpkeepItem{}}
	var err error
	if out.Settings, err = ReadUpkeep(ctx, db); err != nil {
		return out, err
	}
	movies, err := mediafiles.MoviesForUpkeep(ctx, db)
	if err != nil {
		return out, err
	}
	var watched map[int64]bool
	seen := false
	if s.WatchedTitles != nil {
		if watched, seen, err = s.WatchedTitles(ctx); err != nil {
			return out, err
		}
	}
	out.WatchHistory = seen
	cut := s.now().AddDate(0, 0, -out.Settings.NeverWatchedDays)
	for _, m := range movies {
		rule := ""
		switch {
		case m.Files == 0 && !m.KeepUpToDate:
			rule = RuleNoFile
		case m.Files > 0 && seen && !watched[m.TitleID] && addedBefore(m.AddedAt, cut):
			rule = RuleNeverWatched
		}
		if rule == "" {
			continue
		}
		if mode := out.Settings.mode(rule); mode != UpkeepOff {
			out.Items = append(out.Items, UpkeepItem{Rule: rule, Mode: mode, MovieUpkeep: m})
		}
	}
	return out, nil
}

func addedBefore(at string, cut time.Time) bool {
	t, err := time.Parse(time.RFC3339Nano, at)
	return err == nil && t.Before(cut)
}

// UpkeepDone is what applying upkeep did.
type UpkeepDone struct {
	Removed   []int64  `json:"removed"`
	Skipped   []int64  `json:"skipped"`
	Problems  []string `json:"problems,omitempty"`
	HeldBytes int64    `json:"held_bytes"`
}

// ApplyUpkeep takes the named movies out of the Library, their files to the
// holding folder (kept for holding_days). Only movies a rule names right
// now are touched; others, and movies with a protected copy, are skipped.
// With auto set, only the rules in auto mode act (the daily run).
func (s *Service) ApplyUpkeep(ctx context.Context, titleIDs []int64, auto bool, by domain.Actor) (UpkeepDone, error) {
	done := UpkeepDone{Removed: []int64{}, Skipped: []int64{}}
	list, err := s.Upkeep(ctx)
	if err != nil {
		return done, err
	}
	named := map[int64]UpkeepItem{}
	for _, it := range list.Items {
		if !auto || it.Mode == UpkeepAuto {
			named[it.TitleID] = it
		}
	}
	if auto {
		titleIDs = titleIDs[:0]
		for id := range named {
			titleIDs = append(titleIDs, id)
		}
		slices.Sort(titleIDs)
	}
	for _, id := range titleIDs {
		it, ok := named[id]
		if !ok {
			done.Skipped = append(done.Skipped, id)
			continue
		}
		_, err := s.Remove(ctx, id, RemoveChoice{DeleteFiles: true}, true, by)
		if err != nil {
			done.Skipped = append(done.Skipped, id)
			if !errors.Is(err, ErrProtectedCopy) && !errors.Is(err, ErrAlreadyRemoved) {
				done.Problems = append(done.Problems, fmt.Sprintf("%s: %v", it.Name, err))
			}
			continue
		}
		done.Removed = append(done.Removed, id)
		done.HeldBytes += it.Size
	}
	return done, nil
}

func (s *Service) handleUpkeep(ctx context.Context, r *jobs.Run) error {
	done, err := s.ApplyUpkeep(ctx, nil, true, JobActor(r))
	r.SetResult(map[string]any{"removed": done.Removed, "skipped": done.Skipped, "problems": done.Problems, "held_bytes": done.HeldBytes})
	if err == nil && len(done.Problems) > 0 {
		return errors.New(done.Problems[0])
	}
	return err
}
