package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (c CarModel) Get(id int64) (*Car, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT *
		FROM cars
		WHERE id = $1`

	var car Car

	err := c.DB.QueryRow(query, id).Scan(
		&car.ID,
		&car.CreatedAt,
		&car.Name,
		&car.Body,
		&car.BrakeSystem,
		&car.Aspiration,
		&car.Horsepower,
		&car.Mpg,
		&car.Cylinders,
		&car.Acceleration,
		&car.Displacement,
		&car.Origin,
		&car.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &car, nil

}

func (c CarModel) Update(car *Car) error {
	query := `
		UPDATE cars
		SET name = $1, body = $2, brake_system = $3, aspiration = $4, horsepower = $5, mpg = $6, cylinders = $7, acceleration = $8, displacement = $9, origin = $10, version = version + 1
		WHERE id = $11 AND version = $12
		RETURNING version`

	args := []interface{}{
		car.Name,
		car.Body,
		car.BrakeSystem,
		car.Aspiration,
		car.Horsepower,
		car.Mpg,
		car.Cylinders,
		car.Acceleration,
		car.Displacement,
		car.Origin,
		car.ID,
		car.Version,
	}

	err := c.DB.QueryRow(query, args...).Scan(&car.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (c CarModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM cars
		WHERE id = $1`

	result, err := c.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (c CarModel) GetAll(title string, filters Filters) ([]*Car, Metadata, error) {
	// Update the SQL query to include the window function which counts the total
	// (filtered) records.
	// (filtered) records.
	query := fmt.Sprintf(`
SELECT count(*) OVER(), id, created_at, name, body, brake_system, aspiration,horsepower,mpg,cylinders,acceleration,displacement,origin, version
FROM cars
WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
ORDER BY %s %s, id ASC
LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []any{title, filters.limit(), filters.offset()}
	rows, err := c.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	defer rows.Close()
	// Declare a totalRecords variable.
	totalRecords := 0
	cars := []*Car{}
	for rows.Next() {
		var car Car
		err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
			&car.ID,
			&car.CreatedAt,
			&car.Name,
			&car.Body,
			&car.BrakeSystem,
			&car.Aspiration,
			&car.Horsepower,
			&car.Mpg,
			&car.Cylinders,
			&car.Acceleration,
			&car.Displacement,
			&car.Origin,
			&car.Version,
		)
		if err != nil {
			return nil, Metadata{}, err // Update this to return an empty Metadata struct.
		}
		cars = append(cars, &car)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// Include the metadata struct when returning.
	return cars, metadata, nil
}
