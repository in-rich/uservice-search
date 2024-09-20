CREATE TABLE IF NOT EXISTS notes (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    author_id         UUID NOT NULL,
    note_id           UUID NOT NULL UNIQUE,

    content      tsvector,
    /* This field content both the target name and the public identifier */
    target_name       tsvector,

    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

--bun:split

CREATE INDEX IF NOT EXISTS note_search ON notes USING GIN ( (notes.content || notes.target_name) );
