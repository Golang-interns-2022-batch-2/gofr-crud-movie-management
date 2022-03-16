package store

import "github.com/iamkakashi/movie-gofr/internal/model"

func GenerateQuery(m *model.Movie) (str string, fieldValues []interface{}) {
	query := "UPDATE MOVIES SET "

	var values []interface{}

	if m.Name.Valid {
		query += "NAME = ?, "

		values = append(values, m.Name.String)
	}

	if m.Genre.Valid {
		query += "GENRE = ?, "

		values = append(values, m.Genre.String)
	}

	if m.Plot.Valid {
		query += "PLOT = ?, "

		values = append(values, m.Plot.String)
	}

	if m.Rating.Valid {
		query += "RATING = ?, "

		values = append(values, m.Rating.Float64)
	}

	if m.ReleaseDate.Valid {
		query += "RELEASEDATE = ?, "

		values = append(values, m.ReleaseDate.String)
	}

	if m.Released.Valid {
		query += "RELEASED = ?, "

		values = append(values, m.Released.Bool)
	}

	if len(values) == 0 {
		return "", nil
	}

	query += " UPDATEDAT=NOW()  WHERE ID = ? AND DELETEAT IS NULL"

	values = append(values, m.ID)

	return query, values
}
