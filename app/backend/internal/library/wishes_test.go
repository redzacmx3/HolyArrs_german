package library

import (
	"context"
	"testing"

	domain "github.com/HolyArrs/holyarrs/app/backend/internal/core/model"
	"github.com/HolyArrs/holyarrs/app/backend/internal/mediafiles"
)

// TestSeasonAndEpisodeWishes: a season's wish, named by its number, reaches
// only that season's episodes, and an episode's wish lays over it.
func TestSeasonAndEpisodeWishes(t *testing.T) {
	s := open(t)
	ctx := context.Background()
	sh := seedLibShow(t, s, "Wish Show")
	svc := &Service{Store: s}
	var ti mediafiles.Title
	if err := mediafiles.LoadTitle(ctx, s.Read(), sh.Title, &ti); err != nil {
		t.Fatal(err)
	}
	two := 2
	out, err := svc.SaveTitleWish(ctx, ti, mediafiles.PickWish{Season: &two, Wish: mediafiles.Wish{Quality: "720p"}}, true, domain.System)
	if err != nil || !out.Saved || out.Items != 2 {
		t.Fatalf("season wish %+v %v", out, err)
	}
	if _, err := svc.SaveTitleWish(ctx, ti, mediafiles.PickWish{ItemID: sh.Eps["S02E02"].Item, Wish: mediafiles.Wish{Languages: []string{"de"}}}, true, domain.System); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.SaveTitleWish(ctx, ti, mediafiles.PickWish{Season: new(int), Wish: mediafiles.Wish{Quality: "720p"}}, false, domain.System); err != mediafiles.ErrWishPlace {
		t.Fatalf("season 0: %v", err)
	}
	pr := mediafiles.PickRules{Body: mediafiles.GoodEveryday()}
	got := func(ep string) mediafiles.PickRulesBody { return mediafiles.RulesFor(ctx, s.Read(), pr, sh.Eps[ep].Item).Body }
	if b := got("S01E01"); b.GoodEnoughAt != "1080p_web" || len(b.Language.Wanted) != 0 {
		t.Fatalf("season 1 changed: %s %v", b.GoodEnoughAt, b.Language)
	}
	if b := got("S02E01"); b.GoodEnoughAt != "720p_bluray" || len(b.Language.Wanted) != 0 {
		t.Fatalf("S02E01: %s %v", b.GoodEnoughAt, b.Language)
	}
	if b := got("S02E02"); b.GoodEnoughAt != "720p_bluray" || len(b.Language.Wanted) != 1 {
		t.Fatalf("S02E02: %s %v", b.GoodEnoughAt, b.Language)
	}
	ws, err := svc.TitleWishes(ctx, ti)
	if err != nil || len(ws.Wishes) != 2 || ws.Wishes[0].Season == nil || *ws.Wishes[0].Season != 2 {
		t.Fatalf("wishes %+v %v", ws, err)
	}
}
