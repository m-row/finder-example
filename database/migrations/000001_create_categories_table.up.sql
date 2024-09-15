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

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name JSONB NOT NULL,
    depth INTEGER NOT NULL DEFAULT 0,
    is_disabled BOOLEAN NOT NULL DEFAULT FALSE,
    is_featured BOOLEAN NOT NULL DEFAULT FALSE,
    parent_id UUID,
    super_parent_id UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT fk_parent_id FOREIGN KEY (parent_id)
    REFERENCES categories (id) ON DELETE CASCADE,

    CONSTRAINT fk_super_parent_id FOREIGN KEY (super_parent_id)
    REFERENCES categories (id) ON DELETE CASCADE,

    CONSTRAINT no_parent_without_super_parent CHECK (
        NOT (
            parent_id IS NOT NULL
            AND super_parent_id IS NULL
        )
    )
);

--trigger: update_update_at
CREATE TRIGGER app_trigger_update_categories_updated_at
BEFORE UPDATE ON categories FOR EACH ROW
EXECUTE PROCEDURE app_func_update_updated_at();
