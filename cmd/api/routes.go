// Filename: cmd/api/routes.go
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(a.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	//Books
	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", a.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/books", a.requireActivatedUser(a.SearchAndlistBooksHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/books/:id", a.requireActivatedUser(a.displayBookHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/books", a.requireActivatedUser(a.createBookHandler))
	router.HandlerFunc(http.MethodPatch, "/api/v1/books/:id", a.requireActivatedUser(a.updateBookHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/books/:id", a.requireActivatedUser(a.deleteBookHandler))

	router.HandlerFunc(http.MethodPut, "/api/v1/users/activated", a.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/tokens/authentication", a.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodPost, "/api/v1/users", a.registerUserHandler)

	return a.recoverPanic(a.rateLimit(a.authenticate(router)))

}
