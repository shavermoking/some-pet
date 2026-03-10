CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    year INTEGER,
    ISBN VARCHAR(255) UNIQUE,
    outOfStock BOOLEAN NOT NULL DEFAULT false,
    rating INTEGER CHECK(rating >= 0 AND rating <= 10),
    created_at   TIMESTAMP DEFAULT NOW(),
    updated_at   TIMESTAMP DEFAULT NOW()
)