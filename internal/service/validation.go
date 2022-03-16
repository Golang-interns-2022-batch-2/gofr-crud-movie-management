package service

import (
	"github.com/iamkakashi/movie-gofr/internal/model"
)

func validateName(m *model.Movie) bool {
	if (m.Name.Valid && m.Name.String == "") || (m.Genre.Valid && m.Genre.String == "") || (m.Plot.Valid && m.Plot.String == "") {
		return false
	}

	return true
}

func Validation(m *model.Movie) bool {
	if !validateName(m) {
		return false
	}

	if m.Rating.Valid && (m.Rating.Float64 < 0 || m.Rating.Float64 > 5.0) {
		return false
	}

	return true
}
