package mediafiles

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/HolyArrs/holyarrs/app/backend/internal/core/database"
)

// PickWish is a Wish saved for a title (PartID and ItemID 0), a season
// (ItemID 0) or one episode.
type PickWish struct {
	TitleID int64 `json:"title_id"`
	PartID  int64 `json:"part_id,omitempty"`
	ItemID  int64 `json:"item_id,omitempty"`
	// Season is the season's number: read back with a season's wish, and
	// may stand for part_id when one is saved.
	Season *int `json:"season,omitempty"`
	Wish   Wish `json:"wish"`
}

// WishesOfTitle lists every wish saved for a title, the title's own first.
func WishesOfTitle(ctx context.Context, q database.Lister, titleID int64) ([]PickWish, error) {
	rows, err := q.QueryContext(ctx, `SELECT title_id, part_id, item_id, wish,
		(SELECT number FROM parts WHERE id = pick_wishes.part_id) FROM pick_wishes
		WHERE title_id = ? AND (part_id = 0 OR part_id IN (SELECT id FROM parts WHERE title_id = pick_wishes.title_id))
		AND (item_id = 0 OR item_id IN (SELECT id FROM items WHERE title_id = pick_wishes.title_id))
		ORDER BY item_id != 0, part_id, item_id`, titleID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []PickWish
	for rows.Next() {
		var w PickWish
		var body string
		if err := rows.Scan(&w.TitleID, &w.PartID, &w.ItemID, &body, &w.Season); err != nil {
			return nil, err
		}
		if json.Unmarshal([]byte(body), &w.Wish) != nil {
			continue
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

// KindWishKey is the settings key (group what_to_pick) of the wish for every
// movie or every show.
func KindWishKey(kind string) string { return kind + "_wish" }

// KindWish is the wish saved for every movie or every show (kind movie or
// show); none for music.
func KindWish(ctx context.Context, q database.Querier, kind string) Wish {
	var w Wish
	if kind != "movie" && kind != "show" {
		return w
	}
	if body, ok, err := database.SettingFound(ctx, q, "what_to_pick", KindWishKey(kind)); err == nil && ok {
		_ = json.Unmarshal([]byte(body), &w)
	}
	return w
}

// SaveKindWishTx saves the wish for every movie or every show.
func SaveKindWishTx(tx *sql.Tx, kind string, w Wish, by database.Actor) error {
	if kind != "movie" && kind != "show" {
		return ErrBadWish
	}
	body, err := json.Marshal(w)
	if err != nil {
		return err
	}
	return database.PutSetting(tx, "what_to_pick", KindWishKey(kind), string(body), by, time.Now())
}

// WishFor is the wish that holds for one item: the one for every movie or
// show, then the title's, its season's and its own, each laid over the one
// before.
func WishFor(ctx context.Context, q interface {
	database.Lister
	database.Querier
}, itemID int64) (Wish, error) {
	var kind string
	if err := q.QueryRowContext(ctx, "SELECT t.kind FROM items i JOIN titles t ON t.id = i.title_id WHERE i.id = ?", itemID).Scan(&kind); err != nil {
		return Wish{}, err
	}
	out := KindWish(ctx, q, kind)
	rows, err := q.QueryContext(ctx, `SELECT w.wish FROM items i JOIN pick_wishes w ON w.title_id = i.title_id
		AND ((w.part_id = 0 AND w.item_id = 0) OR (w.part_id = COALESCE(i.part_id, -1) AND w.item_id = 0) OR w.item_id = i.id)
		WHERE i.id = ? ORDER BY w.item_id != 0, w.part_id != 0`, itemID)
	if err != nil {
		return out, err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var body string
		if err := rows.Scan(&body); err != nil {
			return out, err
		}
		var w Wish
		if json.Unmarshal([]byte(body), &w) == nil {
			out = out.Over(w)
		}
	}
	return out, rows.Err()
}

// RulesFor returns the Pick rules with the item's wish laid on. A failed
// read keeps the rules as they are.
func RulesFor(ctx context.Context, q interface {
	database.Lister
	database.Querier
}, pr PickRules, itemID int64) PickRules {
	w, err := WishFor(ctx, q, itemID)
	if err != nil || w.Empty() {
		return pr
	}
	pr.Body = ApplyWish(pr.Body, w)
	return pr
}

// ErrWishPlace says a wish names a season or episode of another title.
var ErrWishPlace = errors.New("that season or episode is not part of the title")

// SeasonOfWishTx sets PartID from Season when the wish names a season by
// its number.
func SeasonOfWishTx(tx *sql.Tx, w *PickWish) error {
	if w.Season == nil || w.PartID != 0 || w.ItemID != 0 {
		return nil
	}
	err := tx.QueryRow("SELECT id FROM parts WHERE title_id = ? AND kind = 'season' AND number = ?", w.TitleID, *w.Season).Scan(&w.PartID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrWishPlace
	}
	return err
}

// SaveWishTx saves (or, when empty, removes) the wish for a title, season or
// episode.
func SaveWishTx(tx *sql.Tx, w PickWish, by database.Actor) error {
	if err := SeasonOfWishTx(tx, &w); err != nil {
		return err
	}
	if w.PartID != 0 || w.ItemID != 0 {
		var n int
		q := "SELECT count(*) FROM parts WHERE id = ? AND title_id = ?"
		id := w.PartID
		if w.ItemID != 0 {
			q, id = "SELECT count(*) FROM items WHERE id = ? AND title_id = ?", w.ItemID
		}
		if err := tx.QueryRow(q, id, w.TitleID).Scan(&n); err != nil {
			return err
		}
		if n == 0 {
			return ErrWishPlace
		}
	}
	if w.ItemID != 0 {
		w.PartID = 0
	}
	if w.Wish.Empty() {
		_, err := tx.Exec("DELETE FROM pick_wishes WHERE title_id = ? AND part_id = ? AND item_id = ?", w.TitleID, w.PartID, w.ItemID)
		return err
	}
	body, err := json.Marshal(w.Wish)
	if err != nil {
		return err
	}
	who, err := json.Marshal(by)
	if err != nil {
		return err
	}
	now := database.FormatTime(time.Now())
	_, err = tx.Exec(`INSERT INTO pick_wishes (created_at, created_by, updated_at, updated_by, title_id, part_id, item_id, wish)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (title_id, part_id, item_id) DO UPDATE SET wish = excluded.wish, updated_at = excluded.updated_at, updated_by = excluded.updated_by`,
		now, string(who), now, string(who), w.TitleID, w.PartID, w.ItemID, string(body))
	return err
}

// RulesInUse is a Pick rules set and how many titles of one kind use it.
type RulesInUse struct {
	PickRules
	Titles int
}

// PickRulesOfKind lists the Pick rules sets the Versions of every movie or
// every show use, with how many titles use each.
func PickRulesOfKind(ctx context.Context, q interface {
	database.Lister
	database.Querier
}, kind string) ([]RulesInUse, error) {
	rows, err := q.QueryContext(ctx, `SELECT v.pick_rules_id, COUNT(DISTINCT v.title_id) FROM versions v
		JOIN titles t ON t.id = v.title_id WHERE t.kind = ? AND t.state != 'removed' GROUP BY v.pick_rules_id ORDER BY v.pick_rules_id`, kind)
	if err != nil {
		return nil, err
	}
	type use struct{ id, n int64 }
	var uses []use
	for rows.Next() {
		var u use
		if err := rows.Scan(&u.id, &u.n); err != nil {
			_ = rows.Close()
			return nil, err
		}
		uses = append(uses, u)
	}
	_ = rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var out []RulesInUse
	for _, u := range uses {
		var r RulesInUse
		if err := LoadPickRules(ctx, q, u.id, &r.PickRules); err != nil {
			continue
		}
		r.Titles = int(u.n)
		out = append(out, r)
	}
	return out, nil
}
