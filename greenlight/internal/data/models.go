package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record (row, entry) not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Create a Models struct which wraps the MovieModel
// kind of enveloping
type Models struct {
	Permissions PermissionModel // Add a new Permissions field.
	Tokens      TokenModel
	Users       UserModel
	Cars        CarModel
	MotorBikes  MotorbikeModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Permissions: PermissionModel{DB: db}, // Initialize a new PermissionModel instance.
		Tokens:      TokenModel{DB: db},
		Users:       UserModel{DB: db},
		Cars:        CarModel{DB: db},
		MotorBikes:  MotorbikeModel{DB: db},
	}
}
