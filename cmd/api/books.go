package main

import (
	"errors"
	"fmt"
	"net/http"

	// import the data package which contains the definition for Comment
	"github.com/mtechguy/test3/internal/data"
	"github.com/mtechguy/test3/internal/validator"
)

var incomingData struct {
	Content *string `json:"content"`
	Author  *string `json:"author"`
}

func (a *applicationDependencies) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	// create a struct to hold a comment
	// we use struct tags to make the names display in lowercase
	var incomingData struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	}
	// perform the decoding
	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	comment := &data.Comment{
		Content: incomingData.Content,
		Author:  incomingData.Author,
	}
	// Initialize a Validator instance
	v := validator.New()

	data.ValidateComment(v, comment)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors) // implemented later
		return
	}
	err = a.commentModel.Insert(comment)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", incomingData) // delete this
	// Set a Location header. The path to the newly created comment
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/comments/%d", comment.ID))

	data := envelope{
		"comment": comment,
	}
	err = a.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// for now display the result
	fmt.Fprintf(w, "%+v\n", incomingData)
}

func (a *applicationDependencies) displayCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL /v1/comments/:id so that we
	// can use it to query teh comments table. We will
	// implement the readIDParam() function later
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Call Get() to retrieve the comment with the specified id
	comment, err := a.commentModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// display the comment
	data := envelope{
		"comment": comment,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

}

func (a *applicationDependencies) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the URL
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Retrieve the comment from the database
	comment, err := a.commentModel.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			a.notFoundResponse(w, r)
		} else {
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Decode the incoming JSON
	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Update the comment fields based on the incoming data
	if incomingData.Content != nil {
		comment.Content = *incomingData.Content
	}
	if incomingData.Author != nil {
		comment.Author = *incomingData.Author
	}

	// Validate the updated comment
	v := validator.New()
	data.ValidateComment(v, comment)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Perform the update in the database
	err = a.commentModel.Update(comment)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Respond with the updated comment
	data := envelope{
		"comment": comment,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.commentModel.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.IDnotFound(w, r, id) // Pass the ID to the custom message handler
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{
		"message": "comment successfully deleted",
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) listCommentsHandler(w http.ResponseWriter, r *http.Request) {
	// Create a struct to hold the query parameters
	// Later on we will add fields for pagination and sorting (filters)
	var queryParametersData struct {
		Content string
		Author  string
		data.Filters
	}
	// get the query parameters from the URL
	queryParameters := r.URL.Query()
	// Load the query parameters into our struct
	queryParametersData.Content = a.getSingleQueryParameter(
		queryParameters,
		"content",
		"")

	queryParametersData.Author = a.getSingleQueryParameter(
		queryParameters,
		"author",
		"")

	v := validator.New()
	queryParametersData.Filters.Page = a.getSingleIntegerParameter(
		queryParameters, "page", 1, v)
	queryParametersData.Filters.PageSize = a.getSingleIntegerParameter(
		queryParameters, "page_size", 10, v)

	queryParametersData.Filters.Sort = a.getSingleQueryParameter(
		queryParameters, "sort", "id")

	queryParametersData.Filters.SortSafeList = []string{"id", "author",
		"-id", "-author"}

	// Check if our filters are valid
	data.ValidateFilters(v, queryParametersData.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	comments, metadata, err := a.commentModel.GetAll(
		queryParametersData.Content,
		queryParametersData.Author,
		queryParametersData.Filters,
	)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	data := envelope{
		"comments":  comments,
		"@metadata": metadata,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
