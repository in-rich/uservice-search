CREATE TABLE IF NOT EXISTS reminders (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    author_id         VARCHAR(255) NOT NULL,
    reminder_id       UUID NOT NULL UNIQUE,

    content           tsvector,

    /* This field content both the target name and the public identifier */
    target_name       tsvector,

    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expired_at        TIMESTAMP WITH TIME ZONE NOT NULL
);
