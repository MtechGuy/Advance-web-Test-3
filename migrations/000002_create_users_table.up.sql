-- Create table for Users
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    username text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    version integer NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS bookreviews (
    id bigserial PRIMARY KEY,
    book_id INT DEFAULT 0 REFERENCES books(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,  -- Allow NULL by removing the DEFAULT 0
    rating FLOAT CHECK (rating BETWEEN 1 AND 5),
    review TEXT,
    review_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    version integer NOT NULL DEFAULT 1
);


-- Create the readinglists table if it does not exist
-- Create the readinglists table if it does not exist
CREATE TABLE IF NOT EXISTS readinglists (
    id bigserial PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    created_by INT DEFAULT 0 REFERENCES users(id) ON DELETE CASCADE,
    books INT REFERENCES books(id) ON DELETE CASCADE,  -- Make sure this refers to a valid book id
    status VARCHAR(50) CHECK (status IN ('currently reading', 'completed')),
    version integer NOT NULL DEFAULT 1
);



-- Modify the users table to add reading_lists and reviews columns
ALTER TABLE users
ADD COLUMN reading_lists INT REFERENCES readinglists(id),
ADD COLUMN reviews INT REFERENCES bookreviews(id);


-- Function to calculate and update the average rating
CREATE OR REPLACE FUNCTION automatic_average_rating()
RETURNS TRIGGER AS $$
BEGIN
    -- Update the average rating of the book associated with the new review
    UPDATE books
    SET average_rating = (
        SELECT ROUND(CAST(AVG(rating) AS NUMERIC), 2)
        FROM bookreviews
        WHERE bookreviews.book_id = NEW.book_id
    )
    WHERE id = NEW.book_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger that executes automatic_average_rating() when a review is added, updated, or deleted
CREATE OR REPLACE TRIGGER update_book_rating
AFTER INSERT OR UPDATE OR DELETE ON bookreviews
FOR EACH ROW
EXECUTE FUNCTION automatic_average_rating();