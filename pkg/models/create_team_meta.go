package models

type CreateTeamMeta struct {
	TeamID string `json:"teamID" validate:"required"`
	UserID string `json:"userID" validate:"required"`
}
