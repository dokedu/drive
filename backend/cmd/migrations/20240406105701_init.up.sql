SET statement_timeout = 0;

CREATE TABLE organisations
(
    id                     text        NOT NULL PRIMARY KEY DEFAULT nanoid(),
    name                   text        NOT NULL,
    stripe_customer_id     text        NULL,
    stripe_subscription_id text        NULL,
    created_at             timestamptz NOT NULL             DEFAULT NOW(),
    deleted_at             timestamptz
);

CREATE TABLE files
(
    id              text        NOT NULL PRIMARY KEY DEFAULT nanoid(),
    name            text        NOT NULL,
    mime_type       text        NOT NULL, -- application/vnd.dokedu-apps.folder
    file_size       bigint      NOT NULL             DEFAULT 0,
    parent_id       text        NULL REFERENCES files,
    is_folder       boolean     NOT NULL             DEFAULT FALSE,
    shared_drive    boolean     NOT NULL             DEFAULT FALSE,
    organisation_id text        NOT NULL REFERENCES organisations,
    created_at      timestamptz NOT NULL             DEFAULT NOW(),
    deleted_at      timestamptz NULL
);

CREATE TYPE user_role AS ENUM ('owner', 'admin', 'user');

CREATE TABLE users
(
    id               text        NOT NULL PRIMARY KEY DEFAULT nanoid(),
    role             user_role   NOT NULL,
    organisation_id  text        NOT NULL REFERENCES organisations,
    first_name       text        NOT NULL,
    last_name        text        NOT NULL,
    email            text        NOT NULL UNIQUE,
    password         text        NULL,
    recovery_token   text        NULL,
    recovery_sent_at timestamptz NULL,
    avatar_file_id   text        NULL REFERENCES files,
    created_at       timestamptz                      DEFAULT NOW() NOT NULL,
    deleted_at       timestamptz
);

CREATE TYPE permission_type AS ENUM ('user', 'group', 'domain', 'anyone');
CREATE TYPE permission_role AS ENUM ('viewer', 'manager');

CREATE TABLE file_permissions
(
    id              text,
    file_id         text                      NOT NULL REFERENCES files (id),
    user_id         text                      NULL REFERENCES users (id),
    email_address   text                      NULL REFERENCES users (id),
    permission_type permission_type           NOT NULL,
    permission_role permission_role           NOT NULL,
    created_at      timestamptz DEFAULT NOW() NOT NULL,
    deleted_at      timestamptz
);

CREATE TABLE sessions
(
    id         text        NOT NULL PRIMARY KEY DEFAULT nanoid(),
    user_id    text        NOT NULL REFERENCES users,
    token      text        NOT NULL,
    created_at timestamptz NOT NULL             DEFAULT NOW(),
    deleted_at timestamptz NULL
);
