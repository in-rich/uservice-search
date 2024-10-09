CREATE TABLE IF NOT EXISTS messages (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    team_id           UUID NOT NULL,
    message_id        UUID NOT NULL UNIQUE,

    content           tsvector,
    /* This field content both the target name and the public identifier */
    target_name       tsvector,

    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
