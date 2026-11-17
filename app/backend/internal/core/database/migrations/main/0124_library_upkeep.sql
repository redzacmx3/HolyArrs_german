-- Library upkeep: the daily run of the upkeep rules set to auto. Every rule
-- starts as a suggestion, so the run removes nothing until the owner says so.
INSERT INTO schedules (created_at, created_by, updated_at, updated_by, kind, every_seconds, enabled) VALUES
    (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'), '{"kind":"system"}', strftime('%Y-%m-%dT%H:%M:%fZ', 'now'), '{"kind":"system"}', 'library_upkeep', 86400, 1)
ON CONFLICT (kind) DO NOTHING;
