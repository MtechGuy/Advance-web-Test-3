CREATE TABLE IF NOT EXISTS readinglists (
    id bigserial PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    books INT NOT NULL REFERENCES books(id),
    status VARCHAR(50) CHECK (status IN ('currently reading', 'completed')) NOT NULL,
    version integer NOT NULL DEFAULT 1
);
