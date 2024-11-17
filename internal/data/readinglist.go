package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mtechguy/test3/internal/validator"
)

// each name begins with uppercase so that they are exportable/public

type ReadingList struct {
	ID          int64   `json:"id"`          // Maps to 'id' in SQL
	Name        string  `json:"name"`        // Maps to 'name' in SQL
	Description string  `json:"description"` // Maps to 'description' in SQL
	CreatedBy   int     `json:"created_by"`  // Maps to 'created_by' in SQL
	Books       *int    `json:"books"`       // Nullable column
	Status      *string `json:"status"`      // Maps to 'status' in SQL
	Version     int     `json:"version"`     // Maps to 'version' in SQL
}

type ReadingListModel struct {
	DB *sql.DB
}

func ValidateReadingList(v *validator.Validator, list *ReadingList) {
	// Validate Name
	v.Check(strings.TrimSpace(list.Name) != "", "name", "must be provided")
	v.Check(len(list.Name) <= 100, "name", "must not be more than 100 characters long")

	// Validate Description
	v.Check(strings.TrimSpace(list.Description) != "", "description", "must be provided")
	v.Check(len(list.Description) <= 200, "description", "must not be more than 200 characters long")

	// Validate CreatedBy (Foreign Key)
	v.Check(list.CreatedBy > 0, "created_by", "must be a valid user ID")
}

func (c ReadingListModel) Insert(list *ReadingList) error {
	// Create a context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Check if the CreatedBy user exists
	var createdByExists bool
	err := c.DB.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)", list.CreatedBy).Scan(&createdByExists)
	if err != nil {
		return fmt.Errorf("error checking created_by user: %w", err)
	}
	if !createdByExists {
		return fmt.Errorf("invalid created_by ID: %d does not exist", list.CreatedBy)
	}

	query := `
		INSERT INTO readinglists (name, description, created_by) 
		VALUES ($1, $2, $3) 
		RETURNING id, version;
			 `

	args := []any{list.Name, list.Description, list.CreatedBy}

	err = c.DB.QueryRowContext(ctx, query, args...).Scan(&list.ID, &list.Version)
	if err != nil {
		return fmt.Errorf("error inserting reading list: %w", err)
	}

	return nil
}

// Get a specific Comment from the comments table
func (c ReadingListModel) Get(id int64) (*ReadingList, error) {
	// check if the id is valid
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// the SQL query to be executed against the database table
	query := `
		 SELECT  id, name, description, created_by, books, status, version
		 FROM readinglists
		 WHERE id = $1
	   `
	// declare a variable of type Comment to store the returned comment
	var list ReadingList

	// Set a 3-second context/timer
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, id).Scan(
		&list.ID,
		&list.Name,
		&list.Description,
		&list.CreatedBy,
		&list.Books,
		&list.Status,
		&list.Version,
	)
	// Cont'd on the next slide
	// check for which type of error
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &list, nil
}

func (c ReadingListModel) Update(list *ReadingList) error {
	// The SQL query to be executed against the database table
	// Every time we make an update, we increment the version number
	query := `
			UPDATE readinglists
			SET  name = $1, description = $2, created_by = $3, books = $4, status = $5, version = version + 1
			WHERE id = $6
			RETURNING version
			`

	args := []any{list.Name, list.Description, list.CreatedBy, list.Books, list.Status, list.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&list.Version)

}

func (c ReadingListModel) Delete(id int64) error {

	// check if the id is valid
	if id < 1 {
		return ErrRecordNotFound
	}
	// the SQL query to be executed against the database table
	query := `
        DELETE FROM readinglists
        WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// ExecContext does not return any rows unlike QueryRowContext.
	// It only returns  information about the the query execution
	// such as how many rows were affected
	result, err := c.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// Probably a wrong id was provided or the client is trying to
	// delete an already deleted comment
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil

}

func (c ReadingListModel) GetAll(name string, status string, filters Filters) ([]*ReadingList, Metadata, error) {

	// the SQL query to be executed against the database table
	query := fmt.Sprintf(`
	SELECT COUNT(*) OVER(), id, name, description, created_by, books, status, version
	FROM readinglists
	WHERE (to_tsvector('simple', name) @@
		  plainto_tsquery('simple', $1) OR $1 = '')
	AND (to_tsvector('simple', status) @@
		 plainto_tsquery('simple', $2) OR $2 = '')
	ORDER BY %s %s, id ASC
	LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := c.DB.QueryContext(ctx, query, name, status, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	lists := []*ReadingList{}

	for rows.Next() {
		var list ReadingList
		// Scan the values, including the pointer for 'books'
		err := rows.Scan(&totalRecords,
			&list.ID,
			&list.Name,
			&list.Description,
			&list.CreatedBy,
			&list.Books, // 'Books' will handle NULLs as nil
			&list.Status,
			&list.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		lists = append(lists, &list)
	}

	// Check for errors after the loop
	err = rows.Err()
	if err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return lists, metadata, nil
}

func (c *ReadingListModel) AddBookToList(listID int, bookID *int, status string) error {
	// Validate the provided status
	validStatuses := []string{"currently reading", "completed"}
	isValidStatus := false
	for _, s := range validStatuses {
		if status == s {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return fmt.Errorf("invalid status value: %s", status)
	}

	// Handle the case where bookID is nil (NULL in the database)
	var bookIDToUpdate interface{}
	if bookID != nil {
		bookIDToUpdate = *bookID
	} else {
		bookIDToUpdate = nil
	}

	// Update the `books` and `status` fields in the reading list
	query := `
		UPDATE readinglists
		SET books = $1, status = $2
		WHERE id = $3
	`
	_, err := c.DB.Exec(query, bookIDToUpdate, status, listID)
	if err != nil {
		return fmt.Errorf("error updating book and status in reading list: %v", err)
	}

	return nil
}

func (c *ReadingListModel) RemoveBookFromList(listID int) error {

	// Remove the book and its status from the reading list by setting them to NULL
	query := `
		UPDATE readinglists
		SET books = NULL, status = NULL
		WHERE id = $1
	`
	_, err := c.DB.Exec(query, listID)
	if err != nil {
		return fmt.Errorf("error removing book and status from reading list: %v", err)
	}

	return nil
}
