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

-- Add columns only if they do not exist already
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'users' 
        AND column_name = 'reading_lists'
    ) THEN
        ALTER TABLE users ADD COLUMN reading_lists INT REFERENCES readinglists(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'users' 
        AND column_name = 'reviews'
    ) THEN
        ALTER TABLE users ADD COLUMN reviews INT REFERENCES bookreviews(id);
    END IF;
END $$;
