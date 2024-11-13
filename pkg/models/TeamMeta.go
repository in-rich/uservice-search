package models

type TeamMeta struct {
	TeamID string `json:"teamID" validate:"required,max=255"`
	UserID string `json:"userID" validate:"required,max=255"`
}
