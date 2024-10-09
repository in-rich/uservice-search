package models

type SearchMessages struct {
	TeamID   string `json:"teamID" validate:"required,max=255"`
	Limit    int    `json:"limit" validate:"required,min=0,max=1000"`
	Offset   int    `json:"offset" validate:"min=0"`
	RawQuery string `json:"rawQuery" validate:"max=2048"`
}
