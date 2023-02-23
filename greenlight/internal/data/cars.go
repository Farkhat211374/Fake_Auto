package data

import (
	"database/sql"
	"time"
)

type Car struct {
	ID           int64     `json:"id"` // Unique integer ID for the movie
	CreatedAt    time.Time `json:"-"`  // Timestamp for when the movie is added to our database, "-" directive, hidden in response
	Name         string    `json:"name"`
	Body         string    `json:"body"`
	BrakeSystem  string    `json:"brake_system"`
	Aspiration   string    `json:"aspiration"`
	Horsepower   float64   `json:"horsepower"`
	Mpg          float64   `json:"mpg"`
	Cylinders    int64     `json:"cylinders"`
	Acceleration float64   `json:"acceleration"`
	Displacement float64   `json:"displacement"`
	Origin       string    `json:"origin"`
	Version      int32     `json:"version"` // The version number starts at 1 and will be incremented each
	// time the movie information is updated
}

type CarModel struct {
	DB *sql.DB
}

func (c CarModel) Insert(car *Car) error {
	query := `
		INSERT INTO cars(name,body,brake_system,aspiration,horsepower,mpg,cylinders,acceleration,displacement,origin)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10)
		RETURNING id, created_at, version`

	return c.DB.QueryRow(query, &car.Name, &car.Body, &car.BrakeSystem, &car.Aspiration, &car.Horsepower, &car.Mpg, &car.Cylinders, &car.Acceleration, &car.Displacement, &car.Origin).Scan(&car.ID, &car.CreatedAt, &car.Version)
}

//
//// method for fetching a specific record from the movies table.
//func (c CarModel) Get(id int64) (*Car, error) {
//	if id < 1 {
//		return nil, ErrRecordNotFound
//	}
//
//	query := `
//		SELECT *
//		FROM movies
//		WHERE id = $1`
//
//	var car Car
//
//	err := c.DB.QueryRow(query, id).Scan(
//		&car.ID,
//		&car.CreatedAt,
//		&car.Title,
//		&car.Year,
//		&car.Runtime,
//		pq.Array(&car.Genres),
//		&car.Version,
//	)
//
//	if err != nil {
//		switch {
//		case errors.Is(err, sql.ErrNoRows):
//			return nil, ErrRecordNotFound
//		default:
//			return nil, err
//		}
//	}
//
//	return &car, nil
//
//}
//
//// method for updating a specific record in the movies table.
//func (c CarModel) Update(car *Car) error {
//	query := `
//		UPDATE movies
//		SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
//		WHERE id = $5
//		RETURNING version`
//
//	args := []interface{}{
//		car.Title,
//		car.Year,
//		car.Runtime,
//		pq.Array(car.Genres),
//		car.ID,
//	}
//
//	return c.DB.QueryRow(query, args...).Scan(&car.Version)
//}
//
//// method for deleting a specific record from the movies table.
//func (c CarModel) Delete(id int64) error {
//	if id < 1 {
//		return ErrRecordNotFound
//	}
//	// Construct the SQL query to delete the record.
//	query := `
//		DELETE FROM movies
//		WHERE id = $1`
//
//	result, err := c.DB.Exec(query, id)
//	if err != nil {
//		return err
//	}
//
//	rowsAffected, err := result.RowsAffected()
//	if err != nil {
//		return err
//	}
//
//	if rowsAffected == 0 {
//		return ErrRecordNotFound
//	}
//
//	return nil
//}
