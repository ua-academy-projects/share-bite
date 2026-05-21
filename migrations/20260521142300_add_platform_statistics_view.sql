-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS analytics;

CREATE OR REPLACE VIEW analytics.platform_statistics AS
SELECT
    (SELECT COUNT(*) FROM auth.users) AS total_users,
    (SELECT COUNT(*) FROM auth.user_roles ur JOIN auth.roles r ON r.id = ur.role_id WHERE r.slug = 'admin') AS total_admin_users,
    (SELECT COUNT(*) FROM auth.user_roles ur JOIN auth.roles r ON r.id = ur.role_id WHERE r.slug = 'moderator') AS total_moderator_users,
    (SELECT COUNT(*) FROM auth.user_roles ur JOIN auth.roles r ON r.id = ur.role_id WHERE r.slug = 'user') AS total_regular_users,
    (SELECT COUNT(*) FROM auth.user_roles ur JOIN auth.roles r ON r.id = ur.role_id WHERE r.slug = 'business') AS total_business_role_users,
    (SELECT COUNT(*) FROM auth.users WHERE status = 'active') AS total_active_users,
    (SELECT COUNT(*) FROM auth.users WHERE status = 'muted') AS total_muted_users,
    (SELECT COUNT(*) FROM auth.users WHERE status = 'suspended') AS total_suspended_users,

    (SELECT COUNT(*) FROM guest.customers) AS total_customers,
    (SELECT COUNT(*) FROM guest.posts) AS total_guest_posts,
    (SELECT COUNT(*) FROM guest.comments) AS total_guest_comments,
    (SELECT COUNT(*) FROM guest.post_likes) AS total_guest_post_likes,
    (SELECT COUNT(*) FROM guest.collections) AS total_collections,
    (SELECT COUNT(*) FROM guest.collection_venues) AS total_collection_venues,
    (SELECT COUNT(*) FROM guest.collection_collaborators) AS total_collection_collaborators,
    (SELECT COUNT(*) FROM guest.collection_invitations) AS total_collection_invitations,
    (SELECT COUNT(*) FROM guest.customer_follows) AS total_customer_follows,

    (SELECT COUNT(*) FROM business.org_units) AS total_business_org_units,
    (SELECT COUNT(*) FROM business.posts) AS total_business_posts,
    (SELECT COUNT(*) FROM business.comments) AS total_business_comments,
    (SELECT COUNT(*) FROM business.likes) AS total_business_likes,
    (SELECT COUNT(*) FROM business.boxes) AS total_business_boxes,
    (SELECT COUNT(*) FROM business.box_items) AS total_business_box_items;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS analytics.platform_statistics;
DROP SCHEMA IF EXISTS analytics CASCADE;
-- +goose StatementEnd
