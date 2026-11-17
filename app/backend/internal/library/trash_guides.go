// TRaSH Guides: HolyArrs reads TRaSH Guides (github.com/TRaSH-Guides/Guides, MIT)
// profiles and release formats and turns them into HolyArrs synced traits
// with points. Only the parts HolyArrs can test the same way are taken;
// the rest is reported as skipped, never dropped silently.
package library

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/HolyArrs/holyarrs/app/backend/internal/mediafiles"
)

// TrashProfile is one TRaSH profile, as listed for picking.
type TrashProfile struct {
	ID          string `json:"id"`
	App         string `json:"app"` // sonarr (shows) or radarr (movies)
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`

	ScoreSet        string            `json:"-"`
	MinScore        int               `json:"min_points"`
	CutoffScore     int               `json:"good_enough_points"`
	MinUpgradeScore int               `json:"upgrade_min_points"`
	FormatItems     map[string]string `json:"-"` // name -> trash id
	// Language is the profile's language ("Any", "Original", "German"),
	// Items its qualities and quality groups, Cutoff the item it is good
	// enough at.
	Language string             `json:"-"`
	Items    []TrashProfileItem `json:"-"`
	Cutoff   string             `json:"-"`
}

// TrashProfileItem is one quality, or a group of qualities that rank as
// one (the German profiles' "Merged QPs"), in a TRaSH profile.
type TrashProfileItem struct {
	Name    string   `json:"name"`
	Allowed bool     `json:"allowed"`
	Items   []string `json:"items"`
}

type rawTrashProfile struct {
	TrashID         string            `json:"trash_id"`
	Name            string            `json:"name"`
	Description     string            `json:"trash_description"`
	URL             string            `json:"trash_url"`
	ScoreSet        string            `json:"trash_score_set"`
	MinFormatScore  int               `json:"minFormatScore"`
	CutoffScore     int               `json:"cutoffFormatScore"`
	MinUpgradeScore int               `json:"minUpgradeFormatScore"`
	FormatItems     map[string]string `json:"formatItems"`
	Language        string             `json:"language"`
	Items           []TrashProfileItem `json:"items"`
	Cutoff          string             `json:"cutoff"`
}

// TrashFormat is one TRaSH release format.
type TrashFormat struct {
	ID     string         `json:"trash_id"`
	Name   string         `json:"name"`
	Scores map[string]int `json:"trash_scores"`
	Specs  []TrashSpec    `json:"specifications"`
}

// TrashSpec is one test of a TrashFormat.
type TrashSpec struct {
	Name           string `json:"name"`
	Implementation string `json:"implementation"`
	Negate         bool   `json:"negate"`
	Required       bool   `json:"required"`
	Fields         struct {
		Value json.RawMessage `json:"value"`
		// Size tests give their ends in GB.
		Min float64 `json:"min"`
		Max float64 `json:"max"`
		// ExceptLanguage turns a language test into "any language but this one".
		ExceptLanguage bool `json:"exceptLanguage"`
	} `json:"fields"`
}

// TrashGroupFormat is one format in a TrashGroup.
type TrashGroupFormat struct {
	Name     string `json:"name"`
	ID       string `json:"trash_id"`
	Required bool   `json:"required"`
	Default  bool   `json:"default"`
}

// TrashGroup is a TRaSH format group: formats added to the profiles it names.
type TrashGroup struct {
	Name     string             `json:"name"`
	Default  string             `json:"default"`
	Formats  []TrashGroupFormat `json:"custom_formats"`
	Profiles struct {
		Include map[string]string `json:"include"`
	} `json:"quality_profiles"`
}

// TrashData is everything read for one app.
type TrashData struct {
	App      string
	Profiles []TrashProfile
	Formats  map[string]TrashFormat // by trash id
	Groups   []TrashGroup
}

// TrashGuides fetches TRaSH Guides. Raw and API can point at a test server.
type TrashGuides struct {
	HTTP *http.Client
	Raw  string // raw files base, no trailing "/"
	API  string // GitHub contents API base for listing a folder
	Ref  string

	mu    sync.Mutex
	cache map[string]trashCached
}

type trashCached struct {
	at time.Time
	d  *TrashData
}

// NewTrashGuides reads from GitHub.
func NewTrashGuides() *TrashGuides {
	return &TrashGuides{
		HTTP: &http.Client{Timeout: 30 * time.Second},
		Raw:  "https://raw.githubusercontent.com/TRaSH-Guides/Guides",
		API:  "https://api.github.com/repos/TRaSH-Guides/Guides/contents",
		Ref:  "master",
	}
}

const trashCacheFor = time.Hour

// Load reads one app's profiles, formats and groups (kept for an hour, so
// picking a profile and then applying it reads GitHub once).
func (c *TrashGuides) Load(ctx context.Context, app string) (*TrashData, error) {
	if app != "sonarr" && app != "radarr" {
		return nil, fmt.Errorf("unknown TRaSH app %q", app)
	}
	c.mu.Lock()
	if e, ok := c.cache[app]; ok && time.Since(e.at) < trashCacheFor {
		c.mu.Unlock()
		return e.d, nil
	}
	c.mu.Unlock()
	d := &TrashData{App: app, Formats: map[string]TrashFormat{}}
	dir := "docs/json/" + app
	for _, part := range []string{"quality-profiles", "cf", "cf-groups"} {
		names, err := c.list(ctx, dir+"/"+part)
		if err != nil {
			return nil, err
		}
		bodies := make([][]byte, len(names))
		errs := make([]error, len(names))
		sem := make(chan struct{}, 8)
		var wg sync.WaitGroup
		for i, n := range names {
			wg.Add(1)
			sem <- struct{}{}
			go func() {
				defer func() { <-sem; wg.Done() }()
				bodies[i], errs[i] = c.get(ctx, c.Raw+"/"+c.Ref+"/"+dir+"/"+part+"/"+n)
			}()
		}
		wg.Wait()
		for i, body := range bodies {
			if errs[i] != nil {
				return nil, errs[i]
			}
			switch part {
			case "quality-profiles":
				var r rawTrashProfile
				if json.Unmarshal(body, &r) == nil && r.TrashID != "" {
					d.Profiles = append(d.Profiles, TrashProfile{ID: r.TrashID, App: app, Name: r.Name,
						Description: trashPlain(r.Description), URL: r.URL, ScoreSet: r.ScoreSet,
						MinScore: r.MinFormatScore, CutoffScore: r.CutoffScore, MinUpgradeScore: r.MinUpgradeScore,
						FormatItems: r.FormatItems, Language: r.Language, Items: r.Items, Cutoff: r.Cutoff})
				}
			case "cf":
				var f TrashFormat
				if json.Unmarshal(body, &f) == nil && f.ID != "" {
					d.Formats[f.ID] = f
				}
			case "cf-groups":
				var g TrashGroup
				if json.Unmarshal(body, &g) == nil {
					d.Groups = append(d.Groups, g)
				}
			}
		}
	}
	sort.Slice(d.Profiles, func(i, j int) bool { return d.Profiles[i].Name < d.Profiles[j].Name })
	sort.Slice(d.Groups, func(i, j int) bool { return d.Groups[i].Name < d.Groups[j].Name })
	c.mu.Lock()
	if c.cache == nil {
		c.cache = map[string]trashCached{}
	}
	c.cache[app] = trashCached{time.Now(), d}
	c.mu.Unlock()
	return d, nil
}

func (c *TrashGuides) list(ctx context.Context, path string) ([]string, error) {
	body, err := c.get(ctx, c.API+"/"+path+"?ref="+c.Ref)
	if err != nil {
		return nil, err
	}
	var items []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("TRaSH Guides list %s: %w", path, err)
	}
	var out []string
	for _, it := range items {
		if it.Type == "file" && strings.HasSuffix(it.Name, ".json") {
			out = append(out, it.Name)
		}
	}
	return out, nil
}

func (c *TrashGuides) get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "HolyArrs")
	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("TRaSH Guides could not be reached: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TRaSH Guides answered %s", res.Status)
	}
	return io.ReadAll(io.LimitReader(res.Body, 4<<20))
}

func trashPlain(s string) string {
	s = strings.ReplaceAll(s, "<br>", " ")
	return strings.Join(strings.Fields(s), " ")
}

// Profile finds a profile by id.
func (d *TrashData) Profile(id string) (TrashProfile, bool) {
	for _, p := range d.Profiles {
		if p.ID == id {
			return p, true
		}
	}
	return TrashProfile{}, false
}

// TrashPlanned is one format turned into a synced trait with its points.
type TrashPlanned struct {
	Trait  mediafiles.OwnTrait
	Points int
}

// TrashSkipped is a format HolyArrs left out, and why.
type TrashSkipped struct {
	Name   string   `json:"name"`
	Kinds  []string `json:"kinds,omitempty"` // test kinds HolyArrs cannot run yet
	Reason string   `json:"reason"`
}

// TrashPlan is what applying a profile does.
type TrashPlan struct {
	Profile TrashProfile
	Traits  []TrashPlanned
	Skipped []TrashSkipped
}

// trashSkippedKinds names the test kinds HolyArrs skips, for the app and docs.
var trashSkippedKinds = map[string]string{
	"IndexerFlagSpecification": "indexer flag",
	"ReleaseTypeSpecification": "release type",
	"YearSpecification":        "year",
}

// Sources the old apps number, mapped to HolyArrs sources.
var (
	trashSonarrSources = map[int]string{1: "tv", 2: "tv", 3: "web_dl", 4: "web_rip", 5: "dvd", 6: "bluray", 7: "remux,disc_image"}
	trashRadarrSources = map[int]string{1: "cam", 2: "telesync", 5: "dvd", 6: "tv", 7: "web_dl", 8: "web_rip", 9: "bluray,remux,disc_image"}
)

// TrashMaxPoints is the Pick rules points bound; TRaSH scores past it are cut.
const TrashMaxPoints = mediafiles.MaxPoints

// Build plans applying profile p: its own formats plus the required and
// default formats of the default groups that include it.
func (d *TrashData) Build(p TrashProfile) TrashPlan {
	plan := TrashPlan{Profile: p}
	seen := map[string]bool{}
	var order []string
	add := func(id string) {
		if id != "" && !seen[id] {
			seen[id] = true
			order = append(order, id)
		}
	}
	names := make([]string, 0, len(p.FormatItems))
	for n := range p.FormatItems {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		add(p.FormatItems[n])
	}
	for _, g := range d.Groups {
		if g.Default != "true" {
			continue
		}
		in := false
		for _, id := range g.Profiles.Include {
			in = in || id == p.ID
		}
		if !in {
			continue
		}
		for _, f := range g.Formats {
			if f.Required || f.Default {
				add(f.ID)
			}
		}
	}
	set := p.ScoreSet
	if set == "" {
		set = "default"
	}
	for _, id := range order {
		f, ok := d.Formats[id]
		if !ok {
			plan.Skipped = append(plan.Skipped, TrashSkipped{Name: id, Reason: "listed, but its file is missing"})
			continue
		}
		score, ok := f.Scores[set]
		if !ok {
			score, ok = f.Scores["default"]
		}
		if !ok || score == 0 {
			continue // no points: it would change nothing
		}
		specs, kinds := trashConvert(d.App, f.Specs)
		if len(kinds) > 0 {
			plan.Skipped = append(plan.Skipped, TrashSkipped{Name: f.Name, Kinds: kinds, Reason: "it tests something HolyArrs cannot test yet"})
			continue
		}
		score = max(-TrashMaxPoints, min(TrashMaxPoints, score))
		plan.Traits = append(plan.Traits, TrashPlanned{Points: score, Trait: mediafiles.OwnTrait{
			Name: f.Name, Match: "name", Source: "trash", SourceID: d.App + ":" + f.ID, Specs: specs, Terms: []string{},
		}})
	}
	return plan
}

// trashConvert turns a TRaSH format's tests into trait tests (see ConvertFormat).
func trashConvert(app string, in []TrashSpec) ([]mediafiles.OwnSpec, []string) {
	specs := make([]FormatSpec, 0, len(in))
	for _, s := range in {
		var str string
		var num int
		if json.Unmarshal(s.Fields.Value, &num) == nil {
			str = strconv.Itoa(num)
		} else {
			_ = json.Unmarshal(s.Fields.Value, &str)
		}
		fs := FormatSpec{Name: s.Name, Implementation: s.Implementation, Negate: s.Negate, Required: s.Required,
			Value: str, Min: s.Fields.Min, Max: s.Fields.Max, Except: s.Fields.ExceptLanguage}
		if s.Implementation == "LanguageSpecification" {
			fs.Language = trashLanguage(app, num)
		}
		specs = append(specs, fs)
	}
	return ConvertFormat(app, specs, ConvertOptions{Languages: true})
}

// trashLanguage maps an old app's language id, as TRaSH formats give it,
// to a language test's value: a code, "original" or "any"; "" when
// HolyArrs has none. Ids up to 25 are the same in both apps.
func trashLanguage(app string, id int) string {
	switch id {
	case -2:
		return "original"
	case -1:
		if app == "radarr" {
			return "any"
		}
		return ""
	}
	if id >= 1 && id <= 25 {
		return trashSharedLanguages[id]
	}
	if app == "radarr" {
		return trashRadarrLanguages[id]
	}
	return trashSonarrLanguages[id]
}

var (
	trashSharedLanguages = map[int]string{1: "en", 2: "fr", 3: "es", 4: "de", 5: "it", 6: "da", 7: "nl", 8: "ja", 9: "is",
		10: "zh", 11: "ru", 12: "pl", 13: "vi", 14: "sv", 15: "no", 16: "fi", 17: "tr", 18: "pt", 19: "nl", 20: "el",
		21: "ko", 22: "hu", 23: "he", 24: "lt", 25: "cs"}
	trashRadarrLanguages = map[int]string{26: "hi", 27: "ro", 28: "th", 29: "bg", 30: "pt", 31: "ar", 32: "uk", 33: "fa",
		34: "bn", 35: "sk", 36: "lv", 37: "es", 38: "ca", 39: "hr", 40: "sr", 41: "bs", 42: "et", 43: "ta", 44: "id",
		45: "te", 46: "mk", 47: "sl", 48: "ml", 49: "kn", 50: "sq", 51: "af", 52: "mr", 53: "tl", 54: "ur", 55: "rm",
		56: "mn", 57: "ka"}
	trashSonarrLanguages = map[int]string{26: "ar", 27: "hi", 28: "bg", 29: "ml", 30: "uk", 31: "sk", 32: "th", 33: "pt",
		34: "es", 35: "ro", 36: "lv", 37: "fa", 38: "ca", 39: "hr", 40: "sr", 41: "bs", 42: "et", 43: "ta", 44: "id",
		45: "mk", 46: "sl"}
)
