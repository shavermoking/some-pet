package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"some-pet/internal/models"
	"strings"
)

type Books struct {
	db *sql.DB
}

func NewBooks(db *sql.DB) *Books {
	return &Books{db: db}
}

func (b *Books) Create(ctx context.Context, book models.Book) error {
	_, err := b.db.Exec("INSERT INTO books (title, author, year, rating) values ($1, $2, $3, $4)",
		book.Title, book.Author, book.Year, book.Rating)

	return err
}

func (b *Books) GetByID(ctx context.Context, id int) (models.Book, error) {
	var book models.Book

	err := b.db.QueryRow("SELECT id, title, author, year, outOfStock, rating FROM books WHERE id=$1", id).
		Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.OutOfStock, &book.Rating)
	if errors.Is(err, sql.ErrNoRows) {
		return book, models.ErrBookNotFound
	}

	return book, err
}

func (b *Books) GetAll(ctx context.Context) ([]models.Book, error) {
	rows, err := b.db.Query("SELECT id, title, author, year, outOfStock, rating FROM books")
	if err != nil {
		return nil, err
	}

	books := make([]models.Book, 0)

	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.OutOfStock, &book.Rating); err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, rows.Err()
}

func (b *Books) Delete(ctx context.Context, id int) error {
	_, err := b.db.Exec("DELETE FROM books WHERE id = $1", id)

	return err
}

func (b *Books) Update(ctx context.Context, id int, input models.UpdateBook) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Author != nil {
		setValues = append(setValues, fmt.Sprintf("author=$%d", argId))
		args = append(args, *input.Author)
		argId++
	}

	if input.Year != nil {
		setValues = append(setValues, fmt.Sprintf("year=$%d", argId))
		args = append(args, *input.Year)
		argId++
	}

	if input.Rating != nil {
		setValues = append(setValues, fmt.Sprintf("rating=$%d", argId))
		args = append(args, *input.Rating)
		argId++
	}

	if len(setValues) == 0 {
		return nil
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE books SET %s WHERE id=$%d", setQuery, argId)
	args = append(args, id)

	_, err := b.db.Exec(query, args...)
	return err
}

func (b *Books) MarkOutOfStock(ctx context.Context, id int) error {
	res, err := b.db.Exec("UPDATE books SET outOfStock = true WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return models.ErrBookNotFound
	}

	return nil
}

func (b *Books) GetRecommend(ctx context.Context) ([]models.Book, error) {
	rows, err := b.db.QueryContext(ctx, `
		SELECT id, title, author, year, outOfStock, rating
		FROM books
		ORDER BY 
			CASE WHEN rating IS NOT NULL THEN rating END DESC,
			year DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]models.Book, 0)
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.OutOfStock, &book.Rating); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, rows.Err()
}
