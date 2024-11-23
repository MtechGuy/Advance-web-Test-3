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

	// Books Section
	// =============
	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", a.requireActivatedUser(a.healthcheckHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/books/:bid", a.displayBookHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/books", a.listBooksHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/book/search", a.searchBookHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/books", a.createBookHandler)
	router.HandlerFunc(http.MethodPatch, "/api/v1/books/:bid", a.requireActivatedUser(a.updateBookHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/books/:bid", a.requireActivatedUser(a.deleteBookHandler))

	// Reading Lists Section
	// =====================
	router.HandlerFunc(http.MethodGet, "/api/v1/lists", a.ReadinglistHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/lists/:lid", a.displayReadingListHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/lists", a.createReadingListHandler)
	router.HandlerFunc(http.MethodPatch, "/api/v1/lists/:lid", a.requireActivatedUser(a.updateReadingListHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/lists/:lid", a.requireActivatedUser(a.deleteReadingListHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/lists/:lid/books", a.addReadingListBookHandler)
	router.HandlerFunc(http.MethodDelete, "/api/v1/lists/:lid/books", a.RemoveReadingListBookHandler)

	// Review Section
	// ==============
	router.HandlerFunc(http.MethodPost, "/api/v1/books/:bid/reviews", a.requireActivatedUser(a.createReviewHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/books/:bid/reviews", a.requireActivatedUser(a.bookReviewsHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/books/:bid/reviews/:rid", a.requireActivatedUser(a.displayReviewHandler))
	router.HandlerFunc(http.MethodPatch, "/api/v1/reviews/:rid", a.requireActivatedUser(a.updateReviewHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/reviews/:rid", a.requireActivatedUser(a.deleteReviewHandler))

	// Users
	// =====
	router.HandlerFunc(http.MethodPut, "/api/v1/users/activated", a.activateUserHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/users/:uid", a.requireActivatedUser(a.listUserProfileHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/users/:uid/reviews", a.requireActivatedUser(a.getUserReviewsHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/tokens/authentication", a.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodPost, "/api/v1/users", a.registerUserHandler)

	return a.recoverPanic(a.rateLimit(a.authenticate(router)))
}
