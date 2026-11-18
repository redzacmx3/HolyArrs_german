-- A show that follows a chosen TMDB episode order can read release names by
-- TMDB's own numbering instead (German releases of some kids' shows count
-- that way): '' follows the chosen order, 'tmdb' reads TMDB's numbers.
ALTER TABLE titles ADD COLUMN release_numbers TEXT NOT NULL DEFAULT '';
