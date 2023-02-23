package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Motorbike struct {
	ID           int64     `json:"id"` // Unique integer ID for the movie
	CreatedAt    time.Time `json:"-"`  // Timestamp for when the movie is added to our database, "-" directive, hidden in response
	Name         string    `json:"name"`
	Horsepower   float64   `json:"horsepower"`
	Type         string    `json:"type"`
	Weight       float64   `json:"weight"`
	ThirdPlace   bool      `json:"third_place"`
	Cylinders    int64     `json:"cylinders"`
	Acceleration float64   `json:"acceleration"`
	Displacement float64   `json:"displacement"`
	Origin       string    `json:"origin"`
	Version      int32     `json:"version"` // The version number starts at 1 and will be incremented each
	// time the movie information is updated
}

type MotorbikeModel struct {
	DB *sql.DB
}

func (m MotorbikeModel) Insert(motorbike *Motorbike) error {
	query := `
		INSERT INTO motorbikes(name,horsepower,type,weight,third_place,cylinders,acceleration,displacement,origin)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, version`

	return m.DB.QueryRow(query, &motorbike.Name, &motorbike.Horsepower, &motorbike.Type, &motorbike.Weight, &motorbike.ThirdPlace, &motorbike.Cylinders, &motorbike.Acceleration, &motorbike.Displacement, &motorbike.Origin).Scan(&motorbike.ID, &motorbike.CreatedAt, &motorbike.Version)
}

func (m MotorbikeModel) Get(id int64) (*Motorbike, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT *
		FROM motorbikes
		WHERE id = $1`

	var motorbike Motorbike

	err := m.DB.QueryRow(query, id).Scan(
		&motorbike.ID,
		&motorbike.CreatedAt,
		&motorbike.Name,
		&motorbike.Horsepower,
		&motorbike.Type,
		&motorbike.Weight,
		&motorbike.ThirdPlace,
		&motorbike.Cylinders,
		&motorbike.Acceleration,
		&motorbike.Displacement,
		&motorbike.Origin,
		&motorbike.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &motorbike, nil

}

func (m MotorbikeModel) Update(motorbike *Motorbike) error {
	query := `
		UPDATE motorbikes
		SET name = $1, horsepower = $2, type = $3, weight = $4, third_place = $5, cylinders = $6, acceleration = $7, displacement = $8, origin = $9, version = version + 1
		WHERE id = $10 AND version = $11
		RETURNING version`

	args := []interface{}{
		motorbike.Name,
		motorbike.Horsepower,
		motorbike.Type,
		motorbike.Weight,
		motorbike.ThirdPlace,
		motorbike.Cylinders,
		motorbike.Acceleration,
		motorbike.Displacement,
		motorbike.Origin,
		motorbike.ID,
		motorbike.Version,
	}

	err := m.DB.QueryRow(query, args...).Scan(&motorbike.Version)
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

func (m MotorbikeModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM motorbikes
		WHERE id = $1`

	result, err := m.DB.Exec(query, id)
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

func (m MotorbikeModel) GetAll(name string, filters Filters) ([]*Motorbike, Metadata, error) {

	query := fmt.Sprintf(`
SELECT count(*) OVER(), id, created_at, name, horsepower, type, weight,third_place,cylinders,acceleration,displacement,origin, version
FROM motorbikes
WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
ORDER BY %s %s, id ASC
LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []any{name, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	defer rows.Close()
	// Declare a totalRecords variable.
	totalRecords := 0
	motorbikes := []*Motorbike{}
	for rows.Next() {
		var motorbike Motorbike
		err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
			&motorbike.ID,
			&motorbike.CreatedAt,
			&motorbike.Name,
			&motorbike.Horsepower,
			&motorbike.Type,
			&motorbike.Weight,
			&motorbike.ThirdPlace,
			&motorbike.Cylinders,
			&motorbike.Acceleration,
			&motorbike.Displacement,
			&motorbike.Origin,
			&motorbike.Version,
		)
		if err != nil {
			return nil, Metadata{}, err // Update this to return an empty Metadata struct.
		}
		motorbikes = append(motorbikes, &motorbike)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// Include the metadata struct when returning.
	return motorbikes, metadata, nil
}
