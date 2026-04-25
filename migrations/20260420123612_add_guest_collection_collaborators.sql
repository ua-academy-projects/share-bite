-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS guest.collection_collaborators
(
    collection_id UUID NOT NULL REFERENCES guest.collections(id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES guest.customers(id) ON DELETE CASCADE,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (collection_id, customer_id)
);

CREATE INDEX idx_collection_collaborators_customer_id ON guest.collection_collaborators(customer_id);

CREATE TYPE guest.collection_invitation_status AS ENUM('pending', 'accepted', 'declined');

CREATE TABLE IF NOT EXISTS guest.collection_invitations
(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    collection_id UUID NOT NULL REFERENCES guest.collections(id) ON DELETE CASCADE,

    status guest.collection_invitation_status NOT NULL DEFAULT 'pending',

    inviter_id UUID NOT NULL REFERENCES guest.customers(id) ON DELETE CASCADE,
    invitee_id UUID NOT NULL REFERENCES guest.customers(id) ON DELETE CASCADE,

    last_sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_collection_invitations_unique ON guest.collection_invitations(collection_id, invitee_id);

CREATE INDEX idx_collection_invitations_invitee_status ON guest.collection_invitations(invitee_id, status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS guest.collection_invitations;
DROP TABLE IF EXISTS guest.collection_collaborators;

DROP TYPE IF EXISTS guest.collection_invitation_status;
-- +goose StatementEnd
