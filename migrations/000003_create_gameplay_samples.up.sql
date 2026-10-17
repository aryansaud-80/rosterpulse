CREATE TABLE IF NOT EXISTS gameplay_samples (
    id UUID DEFAULT gen_random_uuid(),
    player_game_profile_id UUID NOT NULL REFERENCES player_game_profiles(id),
    media_asset_id UUID NOT NULL REFERENCES media_assets(id),
    title VARCHAR(150) NOT NULL,
    description TEXT,
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (id)
);