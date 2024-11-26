CREATE TABLE university (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    abbreviated_name VARCHAR(100) NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    domain_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(100) NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(100),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(100)
);


CREATE TABLE "user" (
    id UUID PRIMARY KEY NOT NULL,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_banned BOOLEAN NOT NULL DEFAULT FALSE,
    is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    has_rate_university BOOLEAN NOT NULL DEFAULT FALSE,
    reputation_points INTEGER NOT NULL DEFAULT 0,
    university_id UUID REFERENCES university(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(100) NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(100),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(100)
);

CREATE INDEX IF NOT EXISTS user_username_index ON "user"(username);
CREATE INDEX IF NOT EXISTS user_email_index ON "user"(email);

CREATE TABLE university_rating (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL REFERENCES "user"(id),
    university_id UUID NOT NULL REFERENCES university(id),
    title VARCHAR(100) NOT NULL,
    content VARCHAR(255) NOT NULL,
    university_major VARCHAR(100) NOT NULL,
    facility_rating INTEGER NOT NULL,
    student_organization_rating INTEGER NOT NULL,
    social_environment_rating INTEGER NOT NULL,
    education_quality_rating INTEGER NOT NULL,
    price_to_value_rating INTEGER NOT NULL,
    overall_rating NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(100) NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(100) NOT NULL,
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(100)
);

CREATE TABLE university_rating_point (
  id UUID PRIMARY KEY NOT NULL,
  university_rating_id UUID NOT NULL REFERENCES university_rating(id),
  type VARCHAR(10) NOT NULL,
  content VARCHAR(50) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_by VARCHAR(100) NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_by VARCHAR(100)  NOT NULL
);

CREATE TABLE user_verify_email (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL REFERENCES "user"(id),
    email VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    is_used BOOLEAN NOT NULL DEFAULT FALSE,
    expired_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(100) NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(100),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(100)
);

CREATE TABLE user_device (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL REFERENCES "user"(id),
    brand VARCHAR(100),
    type VARCHAR(255),
    model VARCHAR(255),
    notification_token VARCHAR(255) NOT NULL,
    is_notification_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(100) NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(100),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(100)
);

CREATE TABLE subthread (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(100) NOT NULL,
    followers_count INTEGER NOT NULL DEFAULT 0,
    image_url VARCHAR(255) NOT NULL,
    label_color VARCHAR(100) NOT NULL,
    university_id UUID REFERENCES university(id),
    is_university_subthread BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(100) NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(100),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(100)
);

CREATE INDEX IF NOT EXISTS subthread_name_index ON subthread(name);

CREATE TABLE subthread_follower (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL REFERENCES "user"(id),
    subthread_id UUID NOT NULL REFERENCES subthread(id),
    is_following BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(100) NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(100),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(100)
);

CREATE TABLE thread (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL REFERENCES "user"(id),
    subthread_id UUID NOT NULL REFERENCES subthread(id),
    title VARCHAR(100) NOT NULL,
    content VARCHAR(255) NOT NULL,
    content_summary VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL,
    like_count INTEGER NOT NULL DEFAULT 0,
    dislike_count INTEGER NOT NULL DEFAULT 0,
    comment_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR
);

CREATE INDEX IF NOT EXISTS thread_cursor_index ON thread(created_at,id);

CREATE TABLE thread_activity (
    id UUID PRIMARY KEY NOT NULL,
    actor_id UUID NOT NULL REFERENCES "user"(id),
    actor_email VARCHAR(100) NOT NULL,
    actor_username VARCHAR(100) NOT NULL,
    thread_id UUID NOT NULL REFERENCES thread(id),
    action VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(100),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(100)
);

CREATE TABLE thread_comment (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL REFERENCES "user"(id),
    thread_id UUID NOT NULL REFERENCES thread(id),
    content VARCHAR(255) NOT NULL,
    like_count INTEGER NOT NULL DEFAULT 0,
    dislike_count INTEGER NOT NULL DEFAULT 0,
    reply_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR
);

CREATE TABLE thread_comment_reply (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL REFERENCES "user"(id),
    thread_id UUID NOT NULL REFERENCES thread(id),
    thread_comment_id uuid NOT NULL REFERENCES thread_comment(id),
    content VARCHAR(255) NOT NULL,
    like_count INTEGER NOT NULL DEFAULT 0,
    dislike_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR
);

