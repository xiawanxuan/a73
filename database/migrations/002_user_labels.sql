ALTER TABLE terrain_labels
    ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS is_system BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_terrain_labels_user_id ON terrain_labels(user_id);
CREATE INDEX IF NOT EXISTS idx_terrain_labels_deleted_at ON terrain_labels(deleted_at);

UPDATE terrain_labels SET is_system = TRUE WHERE is_system = FALSE;

DROP INDEX IF EXISTS idx_terrain_labels_name;
CREATE UNIQUE INDEX IF NOT EXISTS idx_terrain_labels_name_user_null
    ON terrain_labels(name) WHERE user_id IS NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_terrain_labels_name_user
    ON terrain_labels(name, user_id) WHERE user_id IS NOT NULL AND deleted_at IS NULL;

CREATE TRIGGER update_terrain_labels_updated_at BEFORE UPDATE ON terrain_labels
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
