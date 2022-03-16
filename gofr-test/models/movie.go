package models

import "time"

type Movie struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Genre       string    `json:"genre"`
	Rating      float64   `json:"rating"`
	ReleaseDate time.Time `json:"releaseDate"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	Plot        string    `json:"plot"`
	Released    bool      `json:"released"`
}
