INSERT INTO games (id, code, name, is_active)
VALUES
    (gen_random_uuid(), 'FREE_FIRE', 'Free Fire', TRUE),
    (gen_random_uuid(), 'PUBG_MOBILE', 'PUBG Mobile', TRUE),
    (gen_random_uuid(), 'MOBILE_LEGENDS', 'Mobile Legends', FALSE)
ON CONFLICT (code) DO NOTHING;