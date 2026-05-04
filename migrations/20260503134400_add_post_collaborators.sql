-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS guest.post_collaborators (
                                                        id BIGSERIAL PRIMARY KEY,

                                                        post_id BIGINT NOT NULL REFERENCES guest.posts(id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES guest.customers(id) ON DELETE CASCADE,

    status VARCHAR(20) NOT NULL DEFAULT 'pending',

    invited_by UUID NOT NULL REFERENCES guest.customers(id),

    expires_at TIMESTAMPTZ NOT NULL,
    responded_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT post_collaborators_unique UNIQUE (post_id, customer_id),

    CONSTRAINT post_collaborators_status_chk
    CHECK (status IN ('pending', 'accepted', 'declined'))
    );
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS guest.post_collaborators;
-- +goose StatementEnd