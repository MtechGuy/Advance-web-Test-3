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

	// Books
	// =====
	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", a.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/books", a.requireActivatedUser(a.SearchAndlistBooksHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/books/:bid", a.requireActivatedUser(a.displayBookHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/books", a.requireActivatedUser(a.createBookHandler))
	router.HandlerFunc(http.MethodPatch, "/api/v1/books/:bid", a.requireActivatedUser(a.updateBookHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/books/:bid", a.requireActivatedUser(a.deleteBookHandler))

	// Reading Lists
	// =============
	router.HandlerFunc(http.MethodGet, "/api/v1/lists", a.requireActivatedUser(a.ReadinglistHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/lists/:lid", a.requireActivatedUser(a.displayReadingListHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/lists", a.requireActivatedUser(a.createReadingListHandler))
	router.HandlerFunc(http.MethodPatch, "/api/v1/lists/:lid", a.requireActivatedUser(a.updateReadingListHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/lists/:lid", a.requireActivatedUser(a.deleteReadingListHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/lists/:lid/books", a.requireActivatedUser(a.addReadingListBookHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/lists/:lid/books", a.requireActivatedUser(a.RemoveReadingListBookHandler))

	// Users
	// =====
	// Define the specific route first
	router.HandlerFunc(http.MethodPut, "/api/v1/users/activated", a.activateUserHandler)
	// Then define the generic route
	router.HandlerFunc(http.MethodGet, "/api/v1/users/:uid", a.requireActivatedUser(a.listUserProfileHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/tokens/authentication", a.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/users", a.registerUserHandler)

	return a.recoverPanic(a.rateLimit(a.authenticate(router)))
}
