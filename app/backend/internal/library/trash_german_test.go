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
