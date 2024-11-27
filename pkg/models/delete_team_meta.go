package models

type DeleteTeamMeta struct {
	TeamID string `json:"teamID" validate:"required,max=255"`
	UserID string `json:"userID" validate:"max=255"`
}
