DROP INDEX IF EXISTS team_and_user_id;
ALTER TABLE IF EXISTS teams_meta DROP CONSTRAINT fk_team_members;
ALTER TABLE IF EXISTS teams_meta DROP COLUMN user_id;

ALTER TABLE IF EXISTS teams_meta ADD COLUMN user_id TEXT NOT NULL;
ALTER TABLE IF EXISTS teams_meta ADD CONSTRAINT fk_team_members UNIQUE (team_id, user_id);
CREATE INDEX team_and_user_id ON teams_meta(team_id,user_id)
