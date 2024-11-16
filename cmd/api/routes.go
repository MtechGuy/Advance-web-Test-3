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

	// GET    /api/v1/books              # List all books with pagination
	// GET    /api/v1/books/{id}         # Get book details
	// POST   /api/v1/books              # Add new book
	// PUT    /api/v1/books/{id}         # Update book details
	// DELETE /api/v1/books/{id}         # Delete book
	// GET    /api/v1/books/search       # Search books by title/author/genre

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comments", a.requireActivatedUser(a.listCommentsHandler))

	router.HandlerFunc(http.MethodPost, "/v1/comments", a.requireActivatedUser(a.createCommentHandler))

	router.HandlerFunc(http.MethodGet, "/v1/comments/:id", a.requireActivatedUser(a.displayCommentHandler))

	router.HandlerFunc(http.MethodPatch, "/v1/comments/:id", a.requireActivatedUser(a.updateCommentHandler))

	router.HandlerFunc(http.MethodDelete, "/v1/comments/:id", a.requireActivatedUser(a.deleteCommentHandler))

	router.HandlerFunc(http.MethodPut, "/v1/users/activated", a.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", a.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", a.registerUserHandler)

	return a.recoverPanic(a.rateLimit(a.authenticate(router)))

}
