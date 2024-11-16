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

-- Create the bookreviews table if it does not exist
CREATE TABLE IF NOT EXISTS bookreviews (
    id bigserial PRIMARY KEY,
    book_id INT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating FLOAT CHECK (rating BETWEEN 1 AND 5),
    review TEXT,
    review_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    version integer NOT NULL DEFAULT 1
);

-- Create the readinglists table if it does not exist
CREATE TABLE IF NOT EXISTS readinglists (
    id bigserial PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    books INT NOT NULL REFERENCES books(id),
    status VARCHAR(50) CHECK (status IN ('currently reading', 'completed')) NOT NULL,
    version integer NOT NULL DEFAULT 1
);

-- Modify the users table to add reading_lists and reviews columns
ALTER TABLE users
ADD COLUMN reading_lists INT REFERENCES readinglists(id),
ADD COLUMN reviews INT REFERENCES bookreviews(id);
