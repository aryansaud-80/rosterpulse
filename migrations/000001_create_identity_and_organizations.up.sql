CREATE TYPE user_role AS ENUM ('PLAYER', 'ORGANIZATION');
CREATE TYPE user_status AS ENUM ('PENDING', 'ACTIVE', 'DISABLED');

CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    normalized_email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'PLAYER',
    status user_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (id)
);

CREATE TYPE media_type AS ENUM ('IMAGE', 'VIDEO');
CREATE TYPE media_status AS ENUM ('PENDING', 'READY', 'FAILED', 'DELETED');

CREATE TABLE IF NOT EXISTS media_assets (
    id UUID DEFAULT gen_random_uuid(),
    uploaded_by_user UUID NOT NULL REFERENCES users(id),
    object_key TEXT NOT NULL,
    media_type media_type NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    file_size_bytes BIGINT NOT NULL CONSTRAINT file_size_bytes_positive CHECK (file_size_bytes > 0),
    status media_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TYPE organization_status AS ENUM ('ACTIVE', 'DISABLED');

CREATE TABLE IF NOT EXISTS organizations (
    id UUID DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    region VARCHAR(100),
    logo_asset_id UUID REFERENCES media_assets(id),
    status organization_status NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (id)
);

CREATE TYPE membership_role AS ENUM ('OWNER', 'MANAGER');
CREATE TYPE membership_status AS ENUM ('ACTIVE', 'REVOKED');

CREATE TABLE IF NOT EXISTS organization_memberships (
    id UUID DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id),
    user_id UUID NOT NULL REFERENCES users(id),
    membership_role membership_role NOT NULL,
    status membership_status NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);