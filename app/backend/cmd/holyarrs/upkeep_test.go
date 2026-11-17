package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"testing"
)

// route: GET /api/v1/upkeep
// route: PUT /api/v1/settings/upkeep
// route: POST /api/v1/upkeep/apply
func TestUpkeepAPI(t *testing.T) {
	t.Parallel()
	url, s, _, titleID, _ := organiseWorld(t)
	b := lbDo(t, "GET", url+"/api/v1/upkeep", "", 200)
	if !strings.Contains(string(b), `"settings"`) {
		t.Fatalf("list %s", b)
	}
	// What the web page sends.
	lbDo(t, "PUT", url+"/api/v1/settings/upkeep", `{"never_watched":"auto","never_watched_days":365,"no_file":"suggest","paused":false}`, 200)
	lbDo(t, "PUT", url+"/api/v1/settings/upkeep", `{"never_watched":"never","never_watched_days":365,"no_file":"suggest","paused":false}`, 400)
	if err := s.Write(context.Background(), func(tx *sql.Tx) error {
		if _, err := tx.Exec("UPDATE files SET state = 'gone' WHERE version_id IN (SELECT id FROM versions WHERE title_id = ?)", titleID); err != nil {
			return err
		}
		_, err := tx.Exec("UPDATE versions SET keep_up_to_date = 0 WHERE title_id = ?", titleID)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	var l struct {
		Items []struct {
			Rule    string `json:"rule"`
			TitleID int64  `json:"title_id"`
		} `json:"items"`
	}
	_ = json.Unmarshal(lbDo(t, "GET", url+"/api/v1/upkeep", "", 200), &l)
	if len(l.Items) != 1 || l.Items[0].Rule != "no_file" {
		t.Fatalf("items %+v", l.Items)
	}
	b = lbDo(t, "POST", url+"/api/v1/upkeep/apply", `{"title_ids":[`+jsonInt(titleID)+`]}`, 200)
	if !strings.Contains(string(b), `"removed":[`+jsonInt(titleID)+`]`) {
		t.Fatalf("apply %s", b)
	}
	// A deleted movie is no longer listed.
	_ = json.Unmarshal(lbDo(t, "GET", url+"/api/v1/upkeep", "", 200), &l)
	if len(l.Items) != 0 {
		t.Fatalf("deleted movie still listed: %+v", l.Items)
	}
}

func jsonInt(n int64) string { b, _ := json.Marshal(n); return string(b) }
