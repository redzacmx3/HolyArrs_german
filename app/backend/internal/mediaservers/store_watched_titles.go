package mediaservers

import (
	"context"

	"github.com/HolyArrs/holyarrs/app/backend/internal/core/database"
)

// WatchedTitles is every Title anyone on a media server watched, and
// whether the watch history holds anything at all (without it, nothing
// can be called unwatched).
func WatchedTitles(ctx context.Context, q database.Lister) (map[int64]bool, bool, error) {
	rows, err := q.QueryContext(ctx, "SELECT DISTINCT title_id FROM watch_state WHERE plays > 0")
	if err != nil {
		return nil, false, err
	}
	defer func() { _ = rows.Close() }()
	out := map[int64]bool{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, false, err
		}
		out[id] = true
	}
	return out, len(out) > 0, rows.Err()
}
