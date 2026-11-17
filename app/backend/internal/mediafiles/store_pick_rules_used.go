package mediafiles

import (
	"context"

	"github.com/HolyArrs/holyarrs/app/backend/internal/core/database"
)

// UsedPickRules lists the Pick rules sets at least one Version uses, by id.
func UsedPickRules(ctx context.Context, q database.Lister) ([]PickRules, error) {
	return database.Select[PickRules](ctx, q, "pick_rules", "id IN (SELECT pick_rules_id FROM versions) ORDER BY id")
}
