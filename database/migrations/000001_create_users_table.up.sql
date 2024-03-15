-- Global updated at update function
CREATE OR REPLACE FUNCTION app_func_update_updated_at()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

CREATE TABLE users (
    id            UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name          TEXT,
    phone         TEXT,
    email         TEXT,
    password_hash BYTEA,
    img           TEXT,
    thumb         TEXT,
    is_disabled   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

ALTER TABLE users ADD CONSTRAINT phone_or_email CHECK ((
    (phone IS NOT NULL)::INTEGER +
    (email IS NOT NULL)::INTEGER
) >= 1);

CREATE UNIQUE INDEX "users_email_unique_nullable"
ON users (email)
WHERE email IS NOT NULL;

CREATE UNIQUE INDEX "users_phone_unique_nullable"
ON users (phone)
WHERE phone IS NOT NULL;

-- trigger: update_update_at
CREATE TRIGGER app_trigger_update_users_updated_at
BEFORE UPDATE ON users FOR EACH ROW
EXECUTE PROCEDURE app_func_update_updated_at();
