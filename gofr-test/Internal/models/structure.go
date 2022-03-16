package models

import (
	"time"
)

type Movie struct {
	ID          int       `json:"Id"`
	Name        string    `json:"Name"`
	Genre       string    `json:"Genre"`
	Rating      float64   `json:"Rating"`
	ReleaseDate time.Time `json:"ReleaseDate"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
	CreatedAt   time.Time `json:"CreatedAt"`
	Plot        string    `json:"Plot"`
	Released    bool      `json:"Released"`
}
type Response struct {
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
	StatusCode int         `json:"statusCode"`
}
