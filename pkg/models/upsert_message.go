package models

type UpsertMessage struct {
	TeamID           string `json:"teamID" validate:"required,max=255"`
	MessageID        string `json:"messageIDID" validate:"required,max=255"`
	Content          string `json:"content" validate:"max=15000"`
	TargetName       string `json:"targetName" validate:"required,max=128"`
	PublicIdentifier string `json:"publicIdentifier" validate:"required,max=128"`
}
