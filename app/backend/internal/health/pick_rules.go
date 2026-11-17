package health

import (
	"context"
	"strconv"
	"strings"

	"github.com/HolyArrs/holyarrs/app/backend/internal/mediafiles"
)

// pickRules checks every Pick rules set a Version uses for settings that
// keep it from picking as its traits ask (mediafiles.AdviseRules): a
// language trait that loses on quality, or good enough at 0 points. Sets
// no Version uses are left out, so unused presets raise nothing.
func (e *Engine) pickRules(ctx context.Context) ([]Result, error) {
	ps, err := mediafiles.UsedPickRules(ctx, e.Store.Read())
	if err != nil {
		// Advice only: a set it cannot read must not stop the other checks.
		return nil, nil
	}
	out := make([]Result, 0, len(ps))
	for _, p := range ps {
		r := Result{Key: "pick_rules:" + strconv.FormatInt(p.ID, 10), Kind: "pick_rules", Subject: p.Name, State: OK}
		if advice := mediafiles.AdviseRules(p.Body); len(advice) > 0 {
			r.State, r.Message = Warning, strings.Join(advice, " ")
		}
		out = append(out, r)
	}
	return out, nil
}
