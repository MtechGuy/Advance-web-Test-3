package main

import (
	"errors"
	"fmt"
	"net/http"

	// import the data package which contains the definition for Comment
	"github.com/mtechguy/test3/internal/data"
	"github.com/mtechguy/test3/internal/validator"
)

// Create a local struct to hold incoming data
var incomingListData struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	CreatedBy   *int    `json:"created_by"`
}

func (a *applicationDependencies) createReadingListHandler(w http.ResponseWriter, r *http.Request) {
	// Create a struct to hold incoming data with the correct field names and JSON tags
	var incomingListData struct {
		Name        string `json:"name"`        // Maps to 'name' in JSON
		Description string `json:"description"` // Maps to 'description' in JSON
		CreatedBy   int    `json:"created_by"`  // Maps to 'created_by' in JSON
	}

	// Perform the decoding of the incoming JSON
	err := a.readJSON(w, r, &incomingListData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	list := &data.ReadingList{
		Name:        incomingListData.Name,
		Description: incomingListData.Description,
		CreatedBy:   incomingListData.CreatedBy,
	}

	// Initialize a Validator instance
	v := validator.New()
	data.ValidateReadingList(v, list)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the reading list into the database
	err = a.readingListModel.Insert(list)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Set a Location header. The path to the newly created reading list
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/v1/lists/%d", list.ID))

	data := envelope{
		"Reading List": list,
	}
	err = a.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) updateReadingListHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the URL
	id, err := a.readIDParam(r, "lid")
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Retrieve the reading list from the database
	list, err := a.readingListModel.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			a.notFoundResponse(w, r)
		} else {
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// // Decode the incoming JSON
	// var incomingListData struct {
	// 	Name        *string `json:"name"`        // Pointer to allow nil
	// 	Description *string `json:"description"` // Pointer to allow nil
	// 	CreatedBy   *int    `json:"created_by"`  // Pointer to allow nil
	// }
	err = a.readJSON(w, r, &incomingListData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Update the fields if new values are provided
	if incomingListData.Name != nil {
		list.Name = *incomingListData.Name
	}
	if incomingListData.Description != nil {
		list.Description = *incomingListData.Description
	}
	if incomingListData.CreatedBy != nil {
		list.CreatedBy = *incomingListData.CreatedBy
	}

	// Validate the updated reading list
	v := validator.New()
	data.ValidateReadingList(v, list)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Perform the update in the database
	err = a.readingListModel.Update(list)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Send the updated reading list as a response
	data := envelope{
		"Reading List": list,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) displayReadingListHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r, "lid")
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	list, err := a.readingListModel.Get(id)
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
		"Reading List": list,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

}

func (a *applicationDependencies) deleteReadingListHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r, "lid")
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.readingListModel.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.BIDnotFound(w, r, id) // Pass the ID to the custom message handler
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{
		"message": "Readling List successfully deleted",
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) ReadinglistHandler(w http.ResponseWriter, r *http.Request) {
	// Create a struct to hold the query parameters
	// Later on we will add fields for pagination and sorting (filters)
	var queryParametersData struct {
		Name   string
		Status string
		data.Filters
	}
	// get the query parameters from the URL
	queryParameters := r.URL.Query()
	// Load the query parameters into our struct
	queryParametersData.Name = a.getSingleQueryParameter(
		queryParameters,
		"name",
		"")

	queryParametersData.Status = a.getSingleQueryParameter(
		queryParameters,
		"status",
		"")

	v := validator.New()
	queryParametersData.Filters.Page = a.getSingleIntegerParameter(
		queryParameters, "page", 1, v)
	queryParametersData.Filters.PageSize = a.getSingleIntegerParameter(
		queryParameters, "page_size", 10, v)

	queryParametersData.Filters.Sort = a.getSingleQueryParameter(
		queryParameters, "sort", "id")

	queryParametersData.Filters.SortSafeList = []string{"id", "name", "status",
		"-id", "-name", "-status"}

	// Check if our filters are valid
	data.ValidateFilters(v, queryParametersData.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	lists, metadata, err := a.readingListModel.GetAll(
		queryParametersData.Name,
		queryParametersData.Status,
		queryParametersData.Filters,
	)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	data := envelope{
		"Reading Lists": lists,
		"@metadata":     metadata,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) addReadingListBookHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the reading list ID from the URL parameter
	id, err := a.readIDParam(r, "lid")
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Modify the incoming request body struct to use sql.NullInt64 for nullable integer
	var incomingListData struct {
		Books  *int    `json:"books"` // Nullable column
		Status *string `json:"status"`
	}

	// Decode the incoming JSON body
	err = a.readJSON(w, r, &incomingListData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Validate that the 'status' field is not nil
	if incomingListData.Status == nil {
		a.badRequestResponse(w, r, fmt.Errorf("status must be provided"))
		return
	}

	// Convert id (int64) to int before passing to AddBookToList
	listID := int(id)

	// Add the book to the reading list
	err = a.readingListModel.AddBookToList(listID, incomingListData.Books, *incomingListData.Status) // Pass the book pointer and status
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Fetch the updated reading list after adding the book
	List, err := a.readingListModel.Get(id)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Set headers for location
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/v1/lists/%d", listID))

	// Send a response with the updated reading list
	data := envelope{
		"Reading List": List, // Use updated list
	}

	err = a.writeJSON(w, http.StatusOK, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) RemoveReadingListBookHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the list ID from the URL parameter
	id, err := a.readIDParam(r, "lid")
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	listID := int(id)
	// Call the RemoveBookFromList method to remove the book from the list
	err = a.readingListModel.RemoveBookFromList(listID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.LIDnotFound(w, r, id) // Pass the ID to the custom message handler
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Return success response
	data := envelope{
		"message": "Book successfully deleted from Reading List",
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
