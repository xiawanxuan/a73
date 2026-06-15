CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "postgis";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(64) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(32) NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS point_clouds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    file_path VARCHAR(512) NOT NULL,
    file_size BIGINT NOT NULL DEFAULT 0,
    point_count BIGINT NOT NULL DEFAULT 0,
    bounds_min_x DOUBLE PRECISION NOT NULL DEFAULT 0,
    bounds_min_y DOUBLE PRECISION NOT NULL DEFAULT 0,
    bounds_min_z DOUBLE PRECISION NOT NULL DEFAULT 0,
    bounds_max_x DOUBLE PRECISION NOT NULL DEFAULT 0,
    bounds_max_y DOUBLE PRECISION NOT NULL DEFAULT 0,
    bounds_max_z DOUBLE PRECISION NOT NULL DEFAULT 0,
    uploaded_by UUID NOT NULL REFERENCES users(id),
    status VARCHAR(32) NOT NULL DEFAULT 'processing',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_point_clouds_uploaded_by ON point_clouds(uploaded_by);
CREATE INDEX IF NOT EXISTS idx_point_clouds_status ON point_clouds(status);
CREATE INDEX IF NOT EXISTS idx_point_clouds_created_at ON point_clouds(created_at DESC);

CREATE TABLE IF NOT EXISTS terrain_labels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(64) NOT NULL UNIQUE,
    color VARCHAR(16) NOT NULL DEFAULT '#1E88E5',
    description TEXT,
    icon VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO terrain_labels (name, color, description, icon) VALUES
    ('海沟', '#1E88E5', '海底深陷的狭长沟槽地形', 'trending_down'),
    ('礁石', '#FF5252', '海中突出的岩石或珊瑚礁', 'terrain'),
    ('水下管线', '#FFB300', '铺设在海底的管道或线缆', 'cable')
ON CONFLICT (name) DO NOTHING;

CREATE TABLE IF NOT EXISTS annotations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    point_cloud_id UUID NOT NULL REFERENCES point_clouds(id) ON DELETE CASCADE,
    label_id UUID NOT NULL REFERENCES terrain_labels(id),
    name VARCHAR(255) NOT NULL,
    polygon_json JSONB NOT NULL,
    bounds_center_x DOUBLE PRECISION,
    bounds_center_y DOUBLE PRECISION,
    bounds_center_z DOUBLE PRECISION,
    creator_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_annotations_point_cloud_id ON annotations(point_cloud_id);
CREATE INDEX IF NOT EXISTS idx_annotations_label_id ON annotations(label_id);
CREATE INDEX IF NOT EXISTS idx_annotations_creator_id ON annotations(creator_id);
CREATE INDEX IF NOT EXISTS idx_annotations_deleted_at ON annotations(deleted_at);

CREATE TABLE IF NOT EXISTS annotation_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    annotation_id UUID NOT NULL REFERENCES annotations(id) ON DELETE CASCADE,
    point_cloud_id UUID NOT NULL REFERENCES point_clouds(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    snapshot_json JSONB NOT NULL,
    operator_id UUID NOT NULL REFERENCES users(id),
    operation VARCHAR(32) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_annotation_version ON annotation_snapshots(annotation_id, version);
CREATE INDEX IF NOT EXISTS idx_snapshots_point_cloud_id ON annotation_snapshots(point_cloud_id);
CREATE INDEX IF NOT EXISTS idx_snapshots_operator_id ON annotation_snapshots(operator_id);
CREATE INDEX IF NOT EXISTS idx_snapshots_created_at ON annotation_snapshots(created_at DESC);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_point_clouds_updated_at BEFORE UPDATE ON point_clouds
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_terrain_labels_updated_at BEFORE UPDATE ON terrain_labels
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_annotations_updated_at BEFORE UPDATE ON annotations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

INSERT INTO users (id, username, password_hash, role) VALUES
    ('00000000-0000-0000-0000-000000000001', 'admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin')
ON CONFLICT (id) DO NOTHING;

INSERT INTO point_clouds (id, name, file_path, file_size, point_count, bounds_min_x, bounds_min_y, bounds_min_z, bounds_max_x, bounds_max_y, bounds_max_z, uploaded_by, status) VALUES
    ('11111111-1111-1111-1111-111111111111', '东海A区块多波束测深_2024Q1', '/data/east_sea_a_2024q1.las', 2147483648, 2500000, -1000, -1000, -500, 1000, 1000, 0, '00000000-0000-0000-0000-000000000001', 'ready'),
    ('22222222-2222-2222-2222-222222222222', '南海B礁盘测深_2024Q2', '/data/south_sea_b_2024q2.las', 1073741824, 1200000, -500, -500, -300, 500, 500, 50, '00000000-0000-0000-0000-000000000001', 'ready'),
    ('33333333-3333-3333-3333-333333333333', '黄海C管道勘测_2024Q3', '/data/yellow_sea_c_2024q3.las', 536870912, 800000, -800, -200, -150, 800, 200, 0, '00000000-0000-0000-0000-000000000001', 'ready')
ON CONFLICT (id) DO NOTHING;
