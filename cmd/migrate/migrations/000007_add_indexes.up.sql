create extension if not exists pg_trgm;

create index if not exists idx_comments_content on comments using gin (content gin_trgm_ops);

create index if not exists idx_posts_title on posts using gin (title gin_trgm_ops);
create index if not exists idx_posts_tags on posts using gin (tags);

create index if not exists idx_users_username on users (username);
create index if not exists idx_posts_user_id on posts (user_id);
create index if not exists idx_comments_post_id on comments (post_id);
