CREATE TYPE profile_status AS ENUM ('ACTIVE', 'HIDDEN', 'DISABLED');

CREATE TABLE IF NOT EXISTS player_profiles (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id),
    display_name VARCHAR(80) NOT NULL,
    bio TEXT,
    region VARCHAR(100),
    profile_image_asset_id UUID REFERENCES media_assets(id),
    status profile_status NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS games (
    id UUID DEFAULT gen_random_uuid(),
    code VARCHAR(40) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS player_game_profiles (
    id UUID DEFAULT gen_random_uuid(),
    player_profile_id UUID NOT NULL REFERENCES player_profiles(id),
    game_id UUID NOT NULL REFERENCES games(id),
    in_game_name VARCHAR(100) NOT NULL,
    in_game_id VARCHAR(100) NOT NULL,
    server_id VARCHAR(50),
    competitive_role VARCHAR(80),
    rank VARCHAR(80),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id),
    CONSTRAINT unique_player_per_game UNIQUE (player_profile_id, game_id),
    CONSTRAINT unique_ingame_id_per_game UNIQUE (game_id, in_game_id)
);

CREATE TABLE IF NOT EXISTS player_team_history (
    id UUID DEFAULT gen_random_uuid(),
    player_game_profile_id UUID NOT NULL REFERENCES player_game_profiles(id),
    organization_id UUID REFERENCES organizations(id),
    team_name VARCHAR(120) NOT NULL,
    start_date DATE,
    end_date DATE,
    is_current BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS player_achievements (
    id UUID DEFAULT gen_random_uuid(),
    player_game_profile_id UUID NOT NULL REFERENCES player_game_profiles(id),
    title VARCHAR(150) NOT NULL,
    tournament VARCHAR(150),
    placement VARCHAR(50),
    achieved_at DATE,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id)
);