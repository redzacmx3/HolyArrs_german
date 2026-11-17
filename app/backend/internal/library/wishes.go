package library

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strconv"
	"time"

	"github.com/HolyArrs/holyarrs/app/backend/internal/core/database"
	domain "github.com/HolyArrs/holyarrs/app/backend/internal/core/model"
	"github.com/HolyArrs/holyarrs/app/backend/internal/mediafiles"
)

// WishRules is what a set of Pick rules asks for, in short: the steps it
// allows, where it is good enough and the languages it wants.
type WishRules struct {
	Name         string   `json:"name"`
	Allowed      []string `json:"allowed"`
	GoodEnoughAt string   `json:"good_enough_at"`
	Languages    []string `json:"languages"`
	OriginalOK   bool     `json:"original_ok"`
}

func wishRules(name string, b mediafiles.PickRulesBody) WishRules {
	r := WishRules{Name: name, GoodEnoughAt: b.GoodEnoughAt, Languages: slices.Clone(b.Language.Wanted), OriginalOK: b.Language.OriginalOK}
	for _, s := range b.Steps {
		if s.Allowed {
			r.Allowed = append(r.Allowed, s.Quality)
		}
	}
	if r.Languages == nil {
		r.Languages = []string{}
	}
	return r
}

// WishVersion is what a wish changes for one Version of a title.
type WishVersion struct {
	VersionID int64     `json:"version_id"`
	Label     string    `json:"label"`
	Before    WishRules `json:"before"`
	After     WishRules `json:"after"`
	// Files are the present files the wish reaches; GoodEnough how many of
	// them are good enough before and after. A file no longer good enough
	// is searched for a better copy when the Version is kept up to date.
	Files            int `json:"files"`
	GoodEnoughBefore int `json:"good_enough_before"`
	GoodEnoughAfter  int `json:"good_enough_after"`
}

// WishOutcome is what a wish for a title, season or episode changes, or
// changed when it was saved.
type WishOutcome struct {
	Items       int           `json:"items"`
	Versions    []WishVersion `json:"versions"`
	SlotChanges int           `json:"slot_changes"`
	Saved       bool          `json:"saved"`
}

var errPreviewOnly = errors.New("preview only")

// TitleWishes are the wishes saved for a title and the one for every title
// of its kind.
type TitleWishes struct {
	Kind   mediafiles.Wish       `json:"kind"`
	Wishes []mediafiles.PickWish `json:"wishes"`
}

// TitleWishes reads the wishes saved for a title.
func (s *Service) TitleWishes(ctx context.Context, t mediafiles.Title) (TitleWishes, error) {
	db := s.Store.Read()
	ws, err := mediafiles.WishesOfTitle(ctx, db, t.ID)
	if ws == nil {
		ws = []mediafiles.PickWish{}
	}
	return TitleWishes{Kind: mediafiles.KindWish(ctx, db, t.Kind), Wishes: ws}, err
}

// SaveTitleWish shows (save false) or saves what a wish for the title, one
// of its seasons or one episode changes. Saving re-checks which files are
// good enough, so a file below the new wish is searched for again.
func (s *Service) SaveTitleWish(ctx context.Context, t mediafiles.Title, pw mediafiles.PickWish, save bool, by domain.Actor) (WishOutcome, error) {
	var out WishOutcome
	if t.Kind != "movie" && t.Kind != "show" {
		return out, mediafiles.ErrBadWish
	}
	if err := pw.Wish.Check(); err != nil {
		return out, err
	}
	pw.TitleID = t.ID
	err := s.Store.Write(ctx, func(tx *sql.Tx) error {
		out = WishOutcome{}
		before, err := measureWish(ctx, tx, t, pw)
		if err != nil {
			return err
		}
		if err := mediafiles.SaveWishTx(tx, pw, by); err != nil {
			return err
		}
		after, err := measureWish(ctx, tx, t, pw)
		if err != nil {
			return err
		}
		for i := range before.vs {
			v := before.vs[i]
			v.After, v.GoodEnoughAfter = after.vs[i].After, after.vs[i].GoodEnoughAfter
			out.Versions = append(out.Versions, v)
		}
		out.Items = before.items
		if !save {
			return errPreviewOnly
		}
		vs, err := mediafiles.VersionsOfTitle(ctx, tx, t.ID)
		if err != nil {
			return err
		}
		for _, v := range vs {
			ch, err := goodEnoughTx(ctx, tx, v, t.KeepUpToDate && v.KeepUpToDate, by)
			if err != nil {
				return err
			}
			out.SlotChanges += len(ch)
		}
		out.Saved = true
		return nil
	})
	if errors.Is(err, errPreviewOnly) {
		err = nil
	}
	if out.Versions == nil {
		out.Versions = []WishVersion{}
	}
	return out, err
}

type wishMeasure struct {
	items int
	vs    []WishVersion
}

// measureWish reads, for each Version of the title, the rules that hold for
// the items the wish reaches and how many of their files are good enough.
// Before and After are both set to the rules now; the caller keeps the half
// it needs.
func measureWish(ctx context.Context, tx *sql.Tx, t mediafiles.Title, pw mediafiles.PickWish) (wishMeasure, error) {
	var m wishMeasure
	where, args := "title_id = ? AND dropped_at IS NULL", []any{t.ID}
	switch {
	case pw.ItemID != 0:
		where, args = "title_id = ? AND id = ?", []any{t.ID, pw.ItemID}
	case pw.PartID != 0:
		where, args = "title_id = ? AND part_id = ? AND dropped_at IS NULL", []any{t.ID, pw.PartID}
	}
	items, err := database.Select[mediafiles.Item](ctx, tx, "items", where+" ORDER BY id", args...)
	if err != nil {
		return m, err
	}
	m.items = len(items)
	vs, err := mediafiles.VersionsOfTitle(ctx, tx, t.ID)
	if err != nil {
		return m, err
	}
	for _, v := range vs {
		var pr mediafiles.PickRules
		if err := database.Get(ctx, tx, "pick_rules", v.PickRulesID, &pr); err != nil {
			return m, err
		}
		wv := WishVersion{VersionID: v.ID, Label: v.Label}
		for i, it := range items {
			rules := mediafiles.RulesFor(ctx, tx, pr, it.ID)
			if i == 0 {
				wv.Before = wishRules(pr.Name, rules.Body)
				wv.After = wv.Before
			}
			sls, err := database.Select[mediafiles.Slot](ctx, tx, "slots", "version_id = ? AND item_id = ? AND file_id IS NOT NULL", v.ID, it.ID)
			if err != nil {
				return m, err
			}
			for _, sl := range sls {
				var f mediafiles.File
				if database.Get(ctx, tx, "files", *sl.FileID, &f) != nil || f.State != "present" {
					continue
				}
				name := f.RelativePath
				if f.ReleaseName != nil && *f.ReleaseName != "" {
					name = *f.ReleaseName
				}
				wv.Files++
				if mediafiles.OwnedReachesGoodEnough(name, f.MediaInfo, rules.Body) && fileMeetsVersionChoices(v, f) {
					wv.GoodEnoughBefore++
				}
			}
		}
		if len(items) == 0 {
			wv.Before = wishRules(pr.Name, pr.Body)
			wv.After = wv.Before
		}
		wv.GoodEnoughAfter = wv.GoodEnoughBefore
		m.vs = append(m.vs, wv)
	}
	return m, nil
}

// KindWishSet is what the wish for every movie or show changes for one Pick
// rules set in use.
type KindWishSet struct {
	PickRulesID int64     `json:"pick_rules_id"`
	Titles      int       `json:"titles"`
	Before      WishRules `json:"before"`
	After       WishRules `json:"after"`
}

// KindWishOutcome is what the wish for every movie or show changes.
type KindWishOutcome struct {
	Kind  string          `json:"kind"`
	Wish  mediafiles.Wish `json:"wish"`
	Sets  []KindWishSet   `json:"sets"`
	Saved bool            `json:"saved"`
}

// KindWish reads the wish for every movie or show and the sets it reaches.
func (s *Service) KindWish(ctx context.Context, kind string) (KindWishOutcome, error) {
	w := mediafiles.KindWish(ctx, s.Store.Read(), kind)
	return s.kindWishSets(ctx, kind, w, w)
}

func (s *Service) kindWishSets(ctx context.Context, kind string, now, next mediafiles.Wish) (KindWishOutcome, error) {
	out := KindWishOutcome{Kind: kind, Wish: next, Sets: []KindWishSet{}}
	if kind != "movie" && kind != "show" {
		return out, mediafiles.ErrBadWish
	}
	sets, err := mediafiles.PickRulesOfKind(ctx, s.Store.Read(), kind)
	if err != nil {
		return out, err
	}
	for _, p := range sets {
		out.Sets = append(out.Sets, KindWishSet{PickRulesID: p.ID, Titles: p.Titles,
			Before: wishRules(p.Name, mediafiles.ApplyWish(p.Body, now)), After: wishRules(p.Name, mediafiles.ApplyWish(p.Body, next))})
	}
	return out, nil
}

// SaveKindWish shows (save false) or saves the wish for every movie or show.
// Saving re-checks which files are good enough for every set it reaches.
func (s *Service) SaveKindWish(ctx context.Context, kind string, w mediafiles.Wish, save bool, by domain.Actor) (KindWishOutcome, error) {
	if err := w.Check(); err != nil {
		return KindWishOutcome{}, err
	}
	out, err := s.kindWishSets(ctx, kind, mediafiles.KindWish(ctx, s.Store.Read(), kind), w)
	if err != nil || !save {
		return out, err
	}
	if err := s.Store.Write(ctx, func(tx *sql.Tx) error { return mediafiles.SaveKindWishTx(tx, kind, w, by) }); err != nil {
		return out, err
	}
	out.Saved = true
	if s.Queue != nil {
		stamp := strconv.FormatInt(time.Now().UnixNano(), 10)
		for _, set := range out.Sets {
			key := KindRecheckPickRules + ":" + strconv.FormatInt(set.PickRulesID, 10) + ":wish:" + stamp
			if _, _, err := s.Queue.Enqueue(ctx, KindRecheckPickRules, &domain.Subject{Type: "pick_rules", ID: set.PickRulesID}, nil, key, by); err != nil {
				return out, err
			}
		}
	}
	return out, nil
}
