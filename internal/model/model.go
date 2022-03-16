package model

import (
	"gopkg.in/guregu/null.v4"
)

type Movie struct {
	ID          int         `json:"id"`
	Name        null.String `json:"name"`
	Genre       null.String `json:"genre"`
	Rating      null.Float  `json:"rating"`
	ReleaseDate null.String `json:"releaseDate"`
	UpdatedAt   null.Time   `json:"updatedAt"`
	CreatedAt   null.Time   `json:"createdAt"`
	Plot        null.String `json:"plot"`
	Released    null.Bool   `json:"released"`
}
