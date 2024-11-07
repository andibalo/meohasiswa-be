
CREATE TABLE "user" (
    id UUID PRIMARY KEY NOT NULL,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_banned BOOLEAN NOT NULL DEFAULT FALSE,
    is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    reputation_points INTEGER NOT NULL DEFAULT 0,
    university_id UUID REFERENCES university(id),
    created_at TIMESTAMPTZ NOT NULL,
    created_by VARCHAR NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR
);

CREATE TABLE user_verify_email (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL REFERENCES "user"(id),
    email VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL,
    expired_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    created_by VARCHAR NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR
);


CREATE TABLE user_device (
    id UUID PRIMARY KEY NOT NULL,
    device_type VARCHAR(100) NOT NULL,
    device_id VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL REFERENCES "user"(id),
    notification_token VARCHAR(255) NOT NULL,
    is_notification_active BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    created_by VARCHAR NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR,
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR
);