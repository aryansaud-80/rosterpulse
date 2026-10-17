CREATE TABLE IF NOT EXISTS watchlist_entries (
    id UUID DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id),
    player_profile_id UUID NOT NULL REFERENCES player_profiles(id),
    saved_by_user_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id),
    CONSTRAINT unique_watchlist_entry UNIQUE (organization_id, player_profile_id)
);

CREATE TYPE trial_status AS ENUM ('PENDING', 'ACCEPTED', 'DECLINED', 'WITHDRAWN', 'CANCELLED', 'COMPLETED');

CREATE TABLE IF NOT EXISTS trial_invitations (
    id UUID DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id),
    player_profile_id UUID NOT NULL REFERENCES player_profiles(id),
    sent_by_user_id UUID NOT NULL REFERENCES users(id),
    status trial_status NOT NULL DEFAULT 'PENDING',
    invitation_note TEXT, 
    scheduled_at TIMESTAMPTZ,
    responded_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id)
);


CREATE UNIQUE INDEX unique_active_trial_invitation 
ON trial_invitations (organization_id, player_profile_id) 
WHERE status IN ('PENDING', 'ACCEPTED');