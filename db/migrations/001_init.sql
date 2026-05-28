CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rooms (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS reservations (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    room_id BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    day DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT reservations_room_slot_unique UNIQUE (room_id, day, start_time, end_time)
);

CREATE INDEX IF NOT EXISTS reservations_room_day_idx ON reservations (room_id, day);
CREATE INDEX IF NOT EXISTS reservations_user_idx ON reservations (user_id);

CREATE TABLE IF NOT EXISTS client_sync (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    config_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    state_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
