package library

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"testing"

	"github.com/HolyArrs/holyarrs/app/backend/internal/core/database"
	"github.com/HolyArrs/holyarrs/app/backend/internal/core/releasename"
	"github.com/HolyArrs/holyarrs/app/backend/internal/mediafiles"
)

// TestTrashGermanFormats runs the German language formats exactly as TRaSH
// Guides ships them (testdata/trash_german, MIT) through the TRaSH converter
// and the judge, and checks each typical German Usenet name fits the formats
// Radarr gives it: German DL for German.DL and German.ML, German for a plain
// German dub, German Subbed for German subtitles.
func TestTrashGermanFormats(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("testdata", "trash_german", "*.json"))
	if err != nil || len(files) == 0 {
		t.Fatalf("testdata: %v %v", files, err)
	}
	var traits []mediafiles.OwnTrait
	byKey := map[string]string{}
	for i, f := range files {
		body, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		var tf TrashFormat
		if err := json.Unmarshal(body, &tf); err != nil {
			t.Fatal(err)
		}
		specs, skipped := trashConvert("radarr", tf.Specs)
		if len(skipped) > 0 {
			t.Fatalf("%s skipped: %v", tf.Name, skipped)
		}
		id := int64(i + 1)
		traits = append(traits, mediafiles.OwnTrait{Meta: database.Meta{ID: id}, Name: tf.Name, Match: "name", Source: "trash", SourceID: "radarr:" + tf.ID, Specs: specs})
		byKey[mediafiles.OwnTraitKey(id)] = tf.Name
	}
	mediafiles.SetOwnTraits(traits)
	t.Cleanup(func() { mediafiles.SetOwnTraits(nil) })

	rules := mediafiles.GoodEveryday()
	rules.Language = mediafiles.LanguageRule{Wanted: []string{"de"}, OriginalOK: true}
	for key := range byKey {
		rules.Traits = append(rules.Traits, mediafiles.TraitPoints{Trait: key, Points: 1})
	}

	fits := func(name, original string) []string {
		p := releasename.Parse(name, releasename.Hints{})
		title := &mediafiles.TitleRef{ID: 1, Name: "Film", Year: 2020, Runtime: 120, OriginalLanguage: original}
		v := mediafiles.Judge(mediafiles.ReleaseFacts{Name: name, Parse: p, Size: 6 << 30, Protocol: "usenet"},
			mediafiles.SlotSet{Match: mediafiles.Match{Title: title, Conf: releasename.Sure}, SlotIDs: []int64{1}, RuntimeMinutes: 120, Out: true}, rules)
		var got []string
		for _, r := range v.Reasons {
			if r.Effect != "add" || r.Step != "points" {
				continue
			}
			if n, ok := byKey[r.Rule[len("trait:"):]]; ok {
				got = append(got, n)
			}
		}
		sort.Strings(got)
		return got
	}

	for _, c := range []struct {
		name, original string
		want           []string
	}{
		{"Film.2020.German.DL.1080p.BluRay.x264-GRP", "en", []string{"German 1080p Booster", "German DL"}},
		{"Film.2020.GERMAN.DL.1080p.WEB.H264-GRP", "en", []string{"German 1080p Booster", "German DL"}},
		{"Film.2020.German.ML.1080p.BluRay.x264-GRP", "en", []string{"German 1080p Booster", "German DL"}},
		{"Film.2020.German.DL.AC3.Dubbed.720p.WEB.H264-GRP", "en", []string{"German DL"}},
		{"Film.2020.German.1080p.WEB.H264-GRP", "en", []string{"German"}},
		{"Film.2020.German.Subbed.1080p.WEB.H264-GRP", "en", []string{"German Subbed"}},
		{"Film.2020.1080p.BluRay.x264-GRP", "en", nil},
		{"Film.2020.1080p.WEB.DL.H264-GRP", "en", nil},
		// No known original language: German DL still holds "the original".
		{"Film.2020.German.DL.1080p.BluRay.x264-GRP", "", []string{"German 1080p Booster", "German DL"}},
		// A German film: German DL is German twice, as the old apps read it.
		{"Film.2020.German.DL.1080p.BluRay.x264-GRP", "de", []string{"German 1080p Booster", "German DL"}},
	} {
		if got := fits(c.name, c.original); !slices.Equal(got, c.want) {
			t.Errorf("%s (original %q): fits %v, want %v", c.name, c.original, got, c.want)
		}
	}
}

// TestTrashGermanProfileRanksGermanFirst applies the shape of TRaSH's
// "[German] HD Bluray + WEB" (testdata/trash_german_profile, MIT): its merged
// 720p/1080p group ranks on points first and it takes any language, so a
// German DL WEB-DL beats an English Blu-ray of the same film, as in Radarr.
func TestTrashGermanProfileRanksGermanFirst(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("testdata", "trash_german_profile", "german-hd-bluray-web.json"))
	if err != nil {
		t.Fatal(err)
	}
	var r rawTrashProfile
	if err := json.Unmarshal(body, &r); err != nil {
		t.Fatal(err)
	}
	prof := TrashProfile{Name: r.Name, Language: r.Language, Items: r.Items, Cutoff: r.Cutoff}
	rules := mediafiles.GoodEveryday()
	trashShape(&rules, prof)
	var allowed []string
	for _, s := range rules.Steps {
		if s.Allowed {
			allowed = append(allowed, s.Quality)
		}
	}
	if !slices.Equal(allowed, []string{"720p_web", "720p_bluray", "1080p_web", "1080p_bluray"}) || rules.GoodEnoughAt != "1080p_bluray" {
		t.Fatalf("steps %v good enough at %s", allowed, rules.GoodEnoughAt)
	}
	if rules.Order != "points_then_quality" || !slices.Contains(rules.Language.Wanted, "de") || !rules.Language.OriginalOK {
		t.Fatalf("order %s language %+v", rules.Order, rules.Language)
	}
	// A sync keeps the languages the owner already chose.
	mine := mediafiles.GoodEveryday()
	mine.Language = mediafiles.LanguageRule{Wanted: []string{"de"}, OriginalOK: true}
	trashShape(&mine, prof)
	if !slices.Equal(mine.Language.Wanted, []string{"de"}) || !mine.Language.OriginalOK {
		t.Fatalf("own choice overwritten: %+v", mine.Language)
	}
	if err := mediafiles.CheckPickRules("German", []string{"movie"}, rules); err != nil {
		t.Fatalf("rules refused: %v", err)
	}

	mediafiles.SetOwnTraits([]mediafiles.OwnTrait{{Meta: database.Meta{ID: 1}, Name: "German DL", Match: "language", Terms: []string{"de", "original"}}})
	t.Cleanup(func() { mediafiles.SetOwnTraits(nil) })
	rules.Traits = []mediafiles.TraitPoints{{Trait: mediafiles.OwnTraitKey(1), Points: 11000}}
	title := &mediafiles.TitleRef{ID: 1, Name: "Film", Year: 2020, Runtime: 120, OriginalLanguage: "en"}
	set := mediafiles.SlotSet{Match: mediafiles.Match{Title: title, Conf: releasename.Sure}, SlotIDs: []int64{1}, RuntimeMinutes: 120, Out: true}
	var rels []mediafiles.ReleaseFacts
	var vs []mediafiles.Verdict
	for _, n := range []string{"Film.2020.1080p.BluRay.x264-ENG", "Film.2020.German.DL.1080p.WEB.H264-GER"} {
		rel := mediafiles.ReleaseFacts{Name: n, Parse: releasename.Parse(n, releasename.Hints{}), Size: 6 << 30, Protocol: "usenet"}
		v := mediafiles.Judge(rel, set, rules)
		if !v.Passed {
			t.Fatalf("%s refused: %+v", n, v.Reasons)
		}
		rels, vs = append(rels, rel), append(vs, v)
	}
	mediafiles.Rank(rels, vs, rules)
	if vs[1].Rank != 1 {
		t.Fatalf("German DL ranked %d (points %d vs %d)", vs[1].Rank, vs[1].Points, vs[0].Points)
	}
}
