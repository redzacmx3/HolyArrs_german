-- A pick wish changes a Version's Pick rules for one movie, show, season or
-- episode: the languages wanted and the resolution aimed for. part_id and
-- item_id are 0 when the wish is for the whole title or season. A wish whose
-- season or episode is gone is left out when read.
CREATE TABLE pick_wishes (
  id INTEGER PRIMARY KEY,
  created_at TEXT NOT NULL,
  created_by TEXT NOT NULL DEFAULT '{}' CHECK(json_valid(created_by)),
  updated_at TEXT NOT NULL,
  updated_by TEXT NOT NULL DEFAULT '{}' CHECK(json_valid(updated_by)),
  title_id INTEGER NOT NULL REFERENCES titles(id) ON DELETE CASCADE,
  part_id INTEGER NOT NULL DEFAULT 0,
  item_id INTEGER NOT NULL DEFAULT 0,
  wish TEXT NOT NULL DEFAULT '{}' CHECK(json_valid(wish)),
  UNIQUE (title_id, part_id, item_id)
);
