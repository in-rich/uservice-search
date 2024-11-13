CREATE TABLE IF NOT EXISTS teams_meta (
id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

team_id           UUID NOT NULL,
user_id           UUID NOT NULL,

CONSTRAINT fk_team_members UNIQUE (team_id, user_id)
);

--bun:split

CREATE INDEX team_and_user_id ON teams_meta(team_id,user_id)
