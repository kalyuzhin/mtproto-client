-- +goose Up
-- +goose StatementBegin
CREATE TABLE tg_users
(
    account_id   BIGINT      NOT NULL,
    user_id      BIGINT      NOT NULL,
    first_name   TEXT        NOT NULL DEFAULT '',
    last_name    TEXT        NOT NULL DEFAULT '',
    username     TEXT        NOT NULL DEFAULT '',
    phone_number TEXT        NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (account_id, user_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tg_users;
-- +goose StatementEnd