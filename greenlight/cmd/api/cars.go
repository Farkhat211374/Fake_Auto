package main

import (
	"errors"
	"fmt"
	"github.com/fara/fakeauto/internal/data"
	"github.com/fara/fakeauto/internal/validator"
	"net/http"
)

func (app *application) createCarHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name         string  `json:"name"`
		Body         string  `json:"body"`
		BrakeSystem  string  `json:"brake_system"`
		Aspiration   string  `json:"aspiration"`
		Horsepower   float64 `json:"horsepower"`
		Mpg          float64 `json:"mpg"`
		Cylinders    int64   `json:"cylinders"`
		Acceleration float64 `json:"acceleration"`
		Displacement float64 `json:"displacement"`
		Origin       string  `json:"origin"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	}

	v := validator.New()

	v.Check(input.Name != "", "name", "must be provided")
	v.Check(len(input.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(input.Body != "", "body", "must be provided")
	v.Check(len(input.Body) <= 50, "body", "must not be more than 50 bytes long")

	v.Check(input.BrakeSystem != "", "brake_system", "must be provided")
	v.Check(len(input.BrakeSystem) <= 50, "brake_system", "must not be more than 50 bytes long")

	v.Check(input.Aspiration != "", "aspiration", "must be provided")
	v.Check(len(input.Aspiration) <= 50, "aspiration", "must not be more than 50 bytes long")

	v.Check(input.Horsepower != 0, "horsepower", "must be provided")
	v.Check(input.Horsepower <= 2000, "horsepower", "must be less than 2000")

	v.Check(input.Mpg != 0, "mpg", "must be provided")

	v.Check(input.Cylinders != 0, "cylinders", "must be provided")
	v.Check(input.Cylinders%2 == 0 && input.Cylinders != 2, "cylinders", "must be 4, 6, 8, 12 etc...")

	v.Check(input.Acceleration != 0, "acceleration", "must be provided")

	v.Check(input.Displacement != 0, "displacement", "must be provided")

	v.Check(input.Origin != "", "origin", "must be provided")
	v.Check(len(input.Origin) <= 50, "origin", "must not be more than 50 bytes long")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)

	car := &data.Car{
		Name:         input.Name,
		Body:         input.Body,
		BrakeSystem:  input.BrakeSystem,
		Aspiration:   input.Aspiration,
		Horsepower:   input.Horsepower,
		Mpg:          input.Mpg,
		Cylinders:    input.Cylinders,
		Acceleration: input.Acceleration,
		Displacement: input.Displacement,
		Origin:       input.Origin,
	}

	err = app.models.Cars.Insert(car)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/cars/%d", car.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"car": car}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showCarHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	car, err := app.models.Cars.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Encode the struct to JSON and send it as the HTTP response.
	// using envelope
	err = app.writeJSON(w, http.StatusOK, envelope{"car": car}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateCarHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Retrieve the movie record as normal.
	car, err := app.models.Cars.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return

	}
	// Use pointers for the Title, Year and Runtime fields.
	var input struct {
		Name         *string  `json:"name"`
		Body         *string  `json:"body"`
		BrakeSystem  *string  `json:"brake_system"`
		Aspiration   *string  `json:"aspiration"`
		Horsepower   *float64 `json:"horsepower"`
		Mpg          *float64 `json:"mpg"`
		Cylinders    *int64   `json:"cylinders"`
		Acceleration *float64 `json:"acceleration"`
		Displacement *float64 `json:"displacement"`
		Origin       *string  `json:"origin"`
	}
	// Decode the JSON as normal.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		car.Name = *input.Name
	}
	// We also do the same for the other fields in the input struct.
	if input.Body != nil {
		car.Body = *input.Body
	}
	if input.BrakeSystem != nil {
		car.BrakeSystem = *input.BrakeSystem
	}
	if input.Aspiration != nil {
		car.Aspiration = *input.Aspiration // Note that we don't need to dereference a slice.
	}
	if input.Horsepower != nil {
		car.Horsepower = *input.Horsepower
	}
	// We also do the same for the other fields in the input struct.
	if input.Mpg != nil {
		car.Mpg = *input.Mpg
	}
	if input.Cylinders != nil {
		car.Cylinders = *input.Cylinders
	}
	if input.Acceleration != nil {
		car.Acceleration = *input.Acceleration // Note that we don't need to dereference a slice.
	}
	if input.Displacement != nil {
		car.Displacement = *input.Displacement
	}
	if input.Origin != nil {
		car.Origin = *input.Origin // Note that we don't need to dereference a slice.
	}

	v := validator.New()

	v.Check(car.Name != "", "name", "must be provided")
	v.Check(len(car.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(car.Body != "", "body", "must be provided")
	v.Check(len(car.Body) <= 50, "body", "must not be more than 50 bytes long")

	v.Check(car.BrakeSystem != "", "brake_system", "must be provided")
	v.Check(len(car.BrakeSystem) <= 50, "brake_system", "must not be more than 50 bytes long")

	v.Check(car.Aspiration != "", "aspiration", "must be provided")
	v.Check(len(car.Aspiration) <= 50, "aspiration", "must not be more than 50 bytes long")

	v.Check(car.Horsepower != 0, "horsepower", "must be provided")
	v.Check(car.Horsepower <= 2000, "horsepower", "must be less than 2000")

	v.Check(car.Mpg != 0, "mpg", "must be provided")

	v.Check(car.Cylinders != 0, "cylinders", "must be provided")
	v.Check(car.Cylinders%2 == 0 && car.Cylinders != 2, "cylinders", "must be 4, 6, 8, 12 etc...")

	v.Check(car.Acceleration != 0, "acceleration", "must be provided")

	v.Check(car.Displacement != 0, "displacement", "must be provided")

	v.Check(car.Origin != "", "origin", "must be provided")
	v.Check(len(car.Origin) <= 50, "origin", "must not be more than 50 bytes long")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)

	err = app.models.Cars.Update(car)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"car": car}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteCarHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Cars.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listCarsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()
	input.Name = app.readString(qs, "name", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "body", "-id", "-name", "-body"}
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Accept the metadata struct as a return value.
	cars, metadata, err := app.models.Cars.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Include the metadata in the response envelope.
	err = app.writeJSON(w, http.StatusOK, envelope{"cars": cars, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
