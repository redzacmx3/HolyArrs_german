package library

import (
	"errors"
	"net/http"

	"github.com/HolyArrs/holyarrs/app/backend/internal/core/webserver"
	"github.com/HolyArrs/holyarrs/app/backend/internal/mediafiles"
)

// wishRoutes serve the language and quality wishes of a title, a season or
// an episode.
func (h *VersionHandlers) wishRoutes(m *http.ServeMux) {
	m.HandleFunc("GET /api/v1/titles/{id}/wishes", h.getWishes)
	m.HandleFunc("PUT /api/v1/titles/{id}/wishes", h.putWish)
	m.HandleFunc("POST /api/v1/titles/{id}/wishes/preview", h.previewWish)
}

func (h *VersionHandlers) getWishes(w http.ResponseWriter, r *http.Request) {
	t, ok := h.TitleID(w, r)
	if !ok {
		return
	}
	ws, err := h.Lib.TitleWishes(r.Context(), t)
	if err != nil {
		webserver.Internal(w, "api", err)
		return
	}
	webserver.WriteJSON(w, http.StatusOK, ws)
}

func (h *VersionHandlers) putWish(w http.ResponseWriter, r *http.Request)     { h.wish(w, r, true) }
func (h *VersionHandlers) previewWish(w http.ResponseWriter, r *http.Request) { h.wish(w, r, false) }

func (h *VersionHandlers) wish(w http.ResponseWriter, r *http.Request, save bool) {
	t, ok := h.TitleID(w, r)
	if !ok {
		return
	}
	var in mediafiles.PickWish
	if !webserver.Decode(w, r, &in, `send {"part_id","item_id","wish":{"languages","original_ok","quality"}}`) {
		return
	}
	out, err := h.Lib.SaveTitleWish(r.Context(), t, in, save, actor(r))
	if err != nil {
		wishFail(w, err)
		return
	}
	if save && h.Publish != nil {
		h.Publish(r, t.ID)
	}
	webserver.WriteJSON(w, http.StatusOK, out)
}

func wishFail(w http.ResponseWriter, err error) {
	if errors.Is(err, mediafiles.ErrBadWish) || errors.Is(err, mediafiles.ErrWishPlace) {
		webserver.Fail(w, http.StatusBadRequest, err.Error())
		return
	}
	webserver.Internal(w, "api", err)
}

// kindWishRoutes serve the wish for every movie and every show.
func (h *OrganiseHandlers) kindWishRoutes(m *http.ServeMux) {
	m.HandleFunc("GET /api/v1/settings/wishes/{kind}", h.getKindWish)
	m.HandleFunc("PUT /api/v1/settings/wishes/{kind}", h.putKindWish)
	m.HandleFunc("POST /api/v1/settings/wishes/{kind}/preview", h.previewKindWish)
}

func (h *OrganiseHandlers) getKindWish(w http.ResponseWriter, r *http.Request) {
	out, err := h.Lib.KindWish(r.Context(), r.PathValue("kind"))
	if errors.Is(err, mediafiles.ErrBadWish) {
		webserver.Fail(w, http.StatusNotFound, "not found")
		return
	}
	if err != nil {
		webserver.Internal(w, "api", err)
		return
	}
	webserver.WriteJSON(w, http.StatusOK, out)
}

func (h *OrganiseHandlers) putKindWish(w http.ResponseWriter, r *http.Request) {
	h.kindWish(w, r, true)
}
func (h *OrganiseHandlers) previewKindWish(w http.ResponseWriter, r *http.Request) {
	h.kindWish(w, r, false)
}

func (h *OrganiseHandlers) kindWish(w http.ResponseWriter, r *http.Request, save bool) {
	var in mediafiles.Wish
	if !webserver.Decode(w, r, &in, `send {"languages","original_ok","quality"}`) {
		return
	}
	out, err := h.Lib.SaveKindWish(r.Context(), r.PathValue("kind"), in, save, actor(r))
	if err != nil {
		wishFail(w, err)
		return
	}
	webserver.WriteJSON(w, http.StatusOK, out)
}
