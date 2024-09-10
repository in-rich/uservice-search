DROP INDEX IF EXISTS note_search;

DROP TRIGGER IF EXISTS search_notes_format ON notes;

--bun:split

DROP FUNCTION IF EXISTS search_notes_format;

DROP TABLE IF EXISTS notes;
