package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"testing"
)

// route: GET /api/v1/titles/{id}/wishes
// route: PUT /api/v1/titles/{id}/wishes
// route: POST /api/v1/titles/{id}/wishes/preview
func TestTitleWishAPI(t *testing.T) {
	t.Parallel()
	url, s, _, titleID, fileID := organiseWorld(t)
	if err := s.Write(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec("UPDATE slots SET file_id = ?, state = 'in_library', good_enough = 1 WHERE version_id IN (SELECT id FROM versions WHERE title_id = ?)", fileID, titleID)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	base := url + "/api/v1/titles/" + jsonInt(titleID) + "/wishes"
	b := lbDo(t, "GET", base, "", 200)
	if !strings.Contains(string(b), `"wishes":[]`) {
		t.Fatalf("wishes %s", b)
	}
	type outcome struct {
		Items    int `json:"items"`
		Versions []struct {
			Files            int `json:"files"`
			GoodEnoughBefore int `json:"good_enough_before"`
			GoodEnoughAfter  int `json:"good_enough_after"`
			After            struct {
				GoodEnoughAt string   `json:"good_enough_at"`
				Languages    []string `json:"languages"`
			} `json:"after"`
		} `json:"versions"`
		SlotChanges int  `json:"slot_changes"`
		Saved       bool `json:"saved"`
	}
	// The 1080p Blu-ray is good enough now; aiming for 4K it is not.
	wish := `{"wish":{"quality":"2160p","languages":["de"],"original_ok":true}}`
	var o outcome
	_ = json.Unmarshal(lbDo(t, "POST", base+"/preview", wish, 200), &o)
	if o.Saved || o.Items != 1 || len(o.Versions) != 1 || o.Versions[0].Files != 1 || o.Versions[0].GoodEnoughBefore != 1 || o.Versions[0].GoodEnoughAfter != 0 ||
		o.Versions[0].After.GoodEnoughAt != "2160p_bluray" || strings.Join(o.Versions[0].After.Languages, ",") != "de" {
		t.Fatalf("preview %+v", o)
	}
	if b := lbDo(t, "GET", base, "", 200); !strings.Contains(string(b), `"wishes":[]`) {
		t.Fatalf("a preview saved: %s", b)
	}
	_ = json.Unmarshal(lbDo(t, "PUT", base, wish, 200), &o)
	if !o.Saved || o.SlotChanges != 1 {
		t.Fatalf("save %+v", o)
	}
	var state string
	if err := s.Read().QueryRow("SELECT s.state FROM slots s JOIN versions v ON v.id = s.version_id WHERE v.title_id = ?", titleID).Scan(&state); err != nil || state != "better_wanted" {
		t.Fatalf("slot %s %v", state, err)
	}
	if b := lbDo(t, "GET", base, "", 200); !strings.Contains(string(b), `"quality":"2160p"`) {
		t.Fatalf("not saved: %s", b)
	}
	lbDo(t, "PUT", base, `{"wish":{"quality":"8k"}}`, 400)
	lbDo(t, "PUT", base, `{"item_id":999999,"wish":{"quality":"720p"}}`, 400)
	// An empty wish takes it away again.
	_ = json.Unmarshal(lbDo(t, "PUT", base, `{"wish":{}}`, 200), &o)
	if o.SlotChanges != 1 {
		t.Fatalf("reset %+v", o)
	}
}

// route: GET /api/v1/settings/wishes/{kind}
// route: PUT /api/v1/settings/wishes/{kind}
// route: POST /api/v1/settings/wishes/{kind}/preview
func TestKindWishAPI(t *testing.T) {
	t.Parallel()
	url, _, _, _, _ := organiseWorld(t)
	base := url + "/api/v1/settings/wishes/movie"
	var o struct {
		Sets []struct {
			Titles int `json:"titles"`
			Before struct {
				GoodEnoughAt string `json:"good_enough_at"`
			} `json:"before"`
			After struct {
				GoodEnoughAt string   `json:"good_enough_at"`
				Allowed      []string `json:"allowed"`
			} `json:"after"`
		} `json:"sets"`
		Saved bool `json:"saved"`
	}
	_ = json.Unmarshal(lbDo(t, "GET", base, "", 200), &o)
	if len(o.Sets) != 1 || o.Sets[0].Titles != 1 {
		t.Fatalf("sets %+v", o)
	}
	_ = json.Unmarshal(lbDo(t, "POST", base+"/preview", `{"quality":"720p"}`, 200), &o)
	if o.Saved || o.Sets[0].After.GoodEnoughAt != "720p_bluray" || strings.Contains(strings.Join(o.Sets[0].After.Allowed, ","), "1080p") {
		t.Fatalf("preview %+v", o)
	}
	_ = json.Unmarshal(lbDo(t, "PUT", base, `{"quality":"720p"}`, 200), &o)
	if !o.Saved {
		t.Fatalf("save %+v", o)
	}
	if b := lbDo(t, "GET", base, "", 200); !strings.Contains(string(b), `"quality":"720p"`) {
		t.Fatalf("not saved: %s", b)
	}
	lbDo(t, "GET", url+"/api/v1/settings/wishes/artist", "", 404)
	lbDo(t, "PUT", base, `{"languages":["German"]}`, 400)
}
