package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/v1/info", app.requireAdminRole(app.createModuleInfoHandler))
	router.HandlerFunc(http.MethodGet, "/v1/info", app.requireActivatedUser(app.getAllModuleInfoHandler))
	router.HandlerFunc(http.MethodGet, "/v1/info/:id", app.requireActivatedUser(app.showModuleInfoHandler))
	router.HandlerFunc(http.MethodPut, "/v1/info/:id", app.requireAdminRole(app.updateModuleInfoHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/info/:id", app.requireAdminRole(app.deleteModuleInfoHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodGet, "/v1/users/get", app.requireActivatedUser(app.getAllUserInfoHandler))
	router.HandlerFunc(http.MethodGet, "/v1/users/get/:id", app.requireActivatedUser(app.getUserInfoHandler))
	router.HandlerFunc(http.MethodPut, "/v1/users/edit/:id", app.requireAdminRole(app.editUserInfoHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/users/delete/:id", app.requireAdminRole(app.deleteUserInfoHandler))

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
