CREATE TABLE IF NOT EXISTS notes (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    author_id         UUID NOT NULL,
    note_id           UUID NOT NULL UNIQUE,

    note_content      tsvector,
    /* This field content both the target name and the public identifier */
    target_name       tsvector,

    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
)

CREATE FUNCTION IF NOT EXISTS search_notes_format()
    RETURN TRIGGER AS $search_notes_format$
BEGIN
    NEW.note_content = setweight(to_tsvector('english', unaccent(NEW.note_content), 'A'));
    NEW.target_name = setweight(to_tsvector('english', unaccent(NEW.target_name), 'A'));
    RETURN NEW;
END;
$search_notes_format$ LANGUAGE plpgsql;

--bun:split

CREATE TRIGGER IF NOT EXISTS search_notes_format
    BEFORE INSERT OR UPDATE ON notes
    FOR EACH ROW
    EXECUTE FUNCTION search_notes_format();

CREATE INDEX IF NOT EXISTS note_search ON notes USING GIN (note_content || target_name);
