package main

import (
	"errors"
	"fmt"
	"github.com/fara/fakeauto/internal/data"
	"github.com/fara/fakeauto/internal/validator"
	"net/http"
)

func (app *application) createMotorbikeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name         string  `json:"name"`
		Horsepower   float64 `json:"horsepower"`
		Type         string  `json:"type"`
		Weight       float64 `json:"weight"`
		ThirdPlace   bool    `json:"third_place"`
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

	v.Check(input.Type != "", "type", "must be provided")

	v.Check(input.Horsepower != 0, "horsepower", "must be provided")
	v.Check(input.Horsepower <= 2000, "horsepower", "must be less than 2000")

	v.Check(input.Weight != 0, "weight", "must be provided")
	v.Check(input.Weight <= 1000, "weight", "must be less than 1000kg")

	v.Check(input.Cylinders != 0, "cylinders", "must be provided")
	v.Check(input.Cylinders%2 == 0, "cylinders", "must be 2, 4 etc...")

	v.Check(input.Acceleration != 0, "acceleration", "must be provided")

	v.Check(input.Displacement != 0, "displacement", "must be provided")

	v.Check(input.Origin != "", "origin", "must be provided")
	v.Check(len(input.Origin) <= 50, "origin", "must not be more than 50 bytes long")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)

	motorbike := &data.Motorbike{
		Name:         input.Name,
		Horsepower:   input.Horsepower,
		Type:         input.Type,
		Weight:       input.Weight,
		ThirdPlace:   input.ThirdPlace,
		Cylinders:    input.Cylinders,
		Acceleration: input.Acceleration,
		Displacement: input.Displacement,
		Origin:       input.Origin,
	}

	err = app.models.MotorBikes.Insert(motorbike)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/motorbikes/%d", motorbike.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"motorbike": motorbike}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showMotorbikeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	motorbike, err := app.models.MotorBikes.Get(id)
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
	err = app.writeJSON(w, http.StatusOK, envelope{"motorbike": motorbike}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateMotorbikeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Retrieve the movie record as normal.
	motorbike, err := app.models.MotorBikes.Get(id)
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
		Horsepower   *float64 `json:"horsepower"`
		Type         *string  `json:"type"`
		Weight       *float64 `json:"weight"`
		ThirdPlace   *bool    `json:"third_place"`
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
		motorbike.Name = *input.Name
	}
	// We also do the same for the other fields in the input struct.
	if input.Horsepower != nil {
		motorbike.Horsepower = *input.Horsepower
	}
	if input.Type != nil {
		motorbike.Type = *input.Type
	}
	if input.Weight != nil {
		motorbike.Weight = *input.Weight // Note that we don't need to dereference a slice.
	}
	// We also do the same for the other fields in the input struct.
	if input.ThirdPlace != nil {
		motorbike.ThirdPlace = *input.ThirdPlace
	}
	if input.Cylinders != nil {
		motorbike.Cylinders = *input.Cylinders
	}
	if input.Acceleration != nil {
		motorbike.Acceleration = *input.Acceleration // Note that we don't need to dereference a slice.
	}
	if input.Displacement != nil {
		motorbike.Displacement = *input.Displacement
	}
	if input.Origin != nil {
		motorbike.Origin = *input.Origin // Note that we don't need to dereference a slice.
	}

	v := validator.New()

	v.Check(motorbike.Name != "", "name", "must be provided")
	v.Check(len(motorbike.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(motorbike.Horsepower != 0, "horsepower", "must be provided")
	v.Check(motorbike.Horsepower <= 2000, "horsepower", "must be less than 2000")

	v.Check(motorbike.Type != "", "type", "must be provided")
	v.Check(len(motorbike.Type) <= 50, "type", "must not be more than 50 bytes long")

	v.Check(motorbike.Weight != 0, "weight", "must be provided")

	v.Check(motorbike.Cylinders != 0, "cylinders", "must be provided")
	v.Check(motorbike.Cylinders%2 == 0, "cylinders", "must be 2, 4 etc...")

	v.Check(motorbike.Acceleration != 0, "acceleration", "must be provided")

	v.Check(motorbike.Displacement != 0, "displacement", "must be provided")

	v.Check(motorbike.Origin != "", "origin", "must be provided")
	v.Check(len(motorbike.Origin) <= 50, "origin", "must not be more than 50 bytes long")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)

	err = app.models.MotorBikes.Update(motorbike)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"motorbike": motorbike}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteMotorbikeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.MotorBikes.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "motorbike successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listMotorbikesHandler(w http.ResponseWriter, r *http.Request) {
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
	input.Filters.SortSafelist = []string{"id", "name", "type", "-id", "-name", "-type"}
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Accept the metadata struct as a return value.
	motorbikes, metadata, err := app.models.MotorBikes.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Include the metadata in the response envelope.
	err = app.writeJSON(w, http.StatusOK, envelope{"motorbikes": motorbikes, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
