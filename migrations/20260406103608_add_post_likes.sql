-- +goose Up
-- +goose StatementBegin
create table if not exists guest.post_likes (
    post_id bigint not null references guest.posts(id) on delete cascade,
    customer_id uuid not null references guest.customers(id) on delete cascade,
    created_at timestamptz not null default now(),
    primary key (post_id, customer_id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
    drop table if exists guest.post_likes;
-- +goose StatementEnd