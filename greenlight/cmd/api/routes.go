package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodGet, "/v1/cars", app.requirePermission("movies:read", app.listCarsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/cars", app.requirePermission("movies:write", app.createCarHandler))
	router.HandlerFunc(http.MethodGet, "/v1/cars/:id", app.requirePermission("movies:read", app.showCarHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/cars/:id", app.requirePermission("movies:write", app.updateCarHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/cars/:id", app.requirePermission("movies:write", app.deleteCarHandler))

	router.HandlerFunc(http.MethodGet, "/v1/motorbikes", app.requirePermission("movies:read", app.listMotorbikesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/motorbikes", app.requirePermission("movies:write", app.createMotorbikeHandler))
	router.HandlerFunc(http.MethodGet, "/v1/motorbikes/:id", app.requirePermission("movies:read", app.showMotorbikeHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/motorbikes/:id", app.requirePermission("movies:write", app.updateMotorbikeHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/motorbikes/:id", app.requirePermission("movies:write", app.deleteMotorbikeHandler))

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
