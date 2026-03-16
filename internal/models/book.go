package models

import (
	"errors"
	"time"
)

var (
	ErrBookNotFound = errors.New("book not found")
)

type Book struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	Year       int    `json:"year"`
	ISBN       string `json:"isbn"`
	OutOfStock bool   `json:"outOfStock"`
	Rating     int    `json:"rating"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CreateBook struct {
	Title  string `json:"title"  binding:"required"`
	Author string `json:"author" binding:"required"`
	Year   int    `json:"year,omitempty"`
	ISBN   string `json:"isbn,omitempty"`
	Rating int    `json:"rating,omitempty"`
}

type UpdateBook struct {
	Title      *string `json:"title"`
	Author     *string `json:"author"`
	Year       *int    `json:"year"`
	ISBN       *string `json:"isbn"`
	OutOfStock *bool   `json:"outOfStock"`
	Rating     *int    `json:"rating"`
}
