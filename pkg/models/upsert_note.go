package models

type UpsertNote struct {
	AuthorID         string `json:"authorID" validate:"required,max=255"`
	NoteID           string `json:"noteID" validate:"required,max=255"`
	Content          string `json:"content" validate:"max=15000"`
	TargetName       string `json:"targetName" validate:"required,max=128"`
	PublicIdentifier string `json:"publicIdentifier" validate:"required,max=128"`
}
