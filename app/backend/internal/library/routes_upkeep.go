package library

import (
	"errors"
	"net/http"

	"github.com/HolyArrs/holyarrs/app/backend/internal/core/database"
	"github.com/HolyArrs/holyarrs/app/backend/internal/core/webserver"
)

// upkeepRoutes serve Library upkeep: what the rules find, their modes,
// and taking the chosen movies out.
func (h *OrganiseHandlers) upkeepRoutes(m *http.ServeMux) {
	m.HandleFunc("GET /api/v1/upkeep", h.getUpkeep)
	m.HandleFunc("PUT /api/v1/settings/upkeep", h.putUpkeep)
	m.HandleFunc("POST /api/v1/upkeep/apply", h.applyUpkeep)
}

func (h *OrganiseHandlers) getUpkeep(w http.ResponseWriter, r *http.Request) {
	l, err := h.Lib.Upkeep(r.Context())
	if err != nil {
		webserver.Internal(w, "api", err)
		return
	}
	webserver.WriteJSON(w, http.StatusOK, l)
}

func (h *OrganiseHandlers) putUpkeep(w http.ResponseWriter, r *http.Request) {
	var u UpkeepSettings
	if !webserver.Decode(w, r, &u, `send {"never_watched","never_watched_days","no_file","paused"}`) {
		return
	}
	if err := h.Lib.SaveUpkeep(r.Context(), u, webserver.ActorOr(r, database.Actor{Kind: "user", ID: 1})); err != nil {
		if errors.Is(err, ErrUpkeepSettings) {
			webserver.Fail(w, http.StatusBadRequest, err.Error())
			return
		}
		webserver.Internal(w, "api", err)
		return
	}
	webserver.WriteJSON(w, http.StatusOK, u)
}

func (h *OrganiseHandlers) applyUpkeep(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TitleIDs []int64 `json:"title_ids"`
	}
	if !webserver.Decode(w, r, &body, `send {"title_ids": [<id>, ...]} from the upkeep list`) {
		return
	}
	if len(body.TitleIDs) == 0 || len(body.TitleIDs) > 1000 {
		webserver.Fail(w, http.StatusBadRequest, "send 1 to 1000 title_ids")
		return
	}
	done, err := h.Lib.ApplyUpkeep(r.Context(), body.TitleIDs, false, webserver.ActorOr(r, database.Actor{Kind: "user", ID: 1}))
	if err != nil {
		webserver.Internal(w, "api", err)
		return
	}
	webserver.WriteJSON(w, http.StatusOK, done)
}
