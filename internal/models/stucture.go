package models

type Movie struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Genre       string  `json:"genre"`
	Rating      float64 `json:"rating"`
	ReleaseDate string  `json:"releaseDate"`
	UpdatedAt   string  `json:"updatedAt"`
	CreatedAt   string  `json:"createdAt"`
	Plot        string  `json:"plot"`
	Released    bool    `json:"released"`
	DeletedAt   string  `json:"deletedAt,omitempty"`
}
