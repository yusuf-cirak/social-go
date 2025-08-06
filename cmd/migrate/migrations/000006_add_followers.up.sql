CREATE TABLE IF NOT EXISTS followers(
    user_id bigint not null,
    follower_id bigint not null,
    created_at timestamp(0) with time zone not null default now(),

    primary key (user_id, follower_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE
)