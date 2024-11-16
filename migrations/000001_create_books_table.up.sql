-- Create table for Books
CREATE TABLE books (
    id bigserial PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    authors TEXT[], -- Array of author names
    isbn VARCHAR(20) UNIQUE,
    publication_date DATE,
    genre VARCHAR(100),
    description TEXT,
    average_rating DECIMAL(3, 2) DEFAULT 0.00, -- Average rating from reviews
    version integer NOT NULL DEFAULT 1
);



