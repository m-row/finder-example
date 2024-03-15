CREATE TABLE roles
(
    id   SMALLSERIAL PRIMARY KEY,
    name TEXT,

    UNIQUE (name)
);

CREATE TABLE permissions
(
    id     SERIAL PRIMARY KEY,
    method TEXT NOT NULL,
    path   TEXT NOT NULL,
    model  TEXT NOT NULL,
    action TEXT NOT NULL,
    scope  TEXT NOT NULL,

    UNIQUE (method, path, model, action, scope)
);

CREATE TABLE role_permissions
(
    role_id       INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,

    PRIMARY KEY (role_id, permission_id),

    CONSTRAINT fk_role_id FOREIGN KEY (role_id)
    REFERENCES roles (id) ON DELETE CASCADE,

    CONSTRAINT fk_permission_id FOREIGN KEY (permission_id)
    REFERENCES permissions (id) ON DELETE CASCADE
);

CREATE TABLE user_roles
(
    user_id UUID NOT NULL,
    role_id INTEGER NOT NULL,

    PRIMARY KEY (user_id, role_id),

    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
    REFERENCES users (id) ON DELETE CASCADE,

    CONSTRAINT fk_role_id FOREIGN KEY (role_id)
    REFERENCES roles (id) ON DELETE CASCADE
);

CREATE TABLE user_permissions
(
    user_id       UUID NOT NULL,
    permission_id INTEGER NOT NULL,

    PRIMARY KEY (user_id, permission_id),

    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
    REFERENCES users (id) ON DELETE CASCADE,

    CONSTRAINT fk_permission_id FOREIGN KEY (permission_id)
    REFERENCES permissions (id) ON DELETE CASCADE
);
