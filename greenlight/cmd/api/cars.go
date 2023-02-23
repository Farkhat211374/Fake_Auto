package main

import (
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
	// // Dump the contents of the input struct in a HTTP response.
	// fmt.Fprintf(w, "%+v\n", input) //+v here is adding the field name of a value // https://pkg.go.dev/fmt
}
