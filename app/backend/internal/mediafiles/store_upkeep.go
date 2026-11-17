package mediafiles

import (
	"context"

	"github.com/HolyArrs/holyarrs/app/backend/internal/core/database"
)

// MovieUpkeep is one movie as the upkeep rules read it: since when it is in
// the Library, its present main files and whether any Version keeps it up
// to date.
type MovieUpkeep struct {
	TitleID      int64  `json:"title_id"`
	Name         string `json:"name"`
	Year         int    `json:"year,omitempty"`
	AddedAt      string `json:"added_at"`
	Size         int64  `json:"size"`
	Files        int    `json:"files"`
	KeepUpToDate bool   `json:"keep_up_to_date"`
}

// MoviesForUpkeep lists every movie with a Version, by sort name.
func MoviesForUpkeep(ctx context.Context, q database.Lister) ([]MovieUpkeep, error) {
	rows, err := q.QueryContext(ctx, `SELECT t.id, t.name, COALESCE(t.year, 0), t.created_at,
		COALESCE(SUM(f.size), 0), COUNT(f.id), MAX(v.keep_up_to_date)
		FROM titles t JOIN versions v ON v.title_id = t.id
		LEFT JOIN files f ON f.version_id = v.id AND f.role = 'main' AND f.state = 'present'
		WHERE t.kind = 'movie' GROUP BY t.id ORDER BY t.sort_name, t.id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []MovieUpkeep
	for rows.Next() {
		var m MovieUpkeep
		var keep int
		if err := rows.Scan(&m.TitleID, &m.Name, &m.Year, &m.AddedAt, &m.Size, &m.Files, &keep); err != nil {
			return nil, err
		}
		m.KeepUpToDate = keep != 0
		out = append(out, m)
	}
	return out, rows.Err()
}
