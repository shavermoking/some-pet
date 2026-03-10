package models

import "errors"

var (
	ErrBookNotFound = errors.New("book not found")
)

type Book struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	Year       int    `json:"year"`
	ISBN       string `json:"ISBN"`
	OutOfStock bool   `json:"outOfStock"`
	Rating     int    `json:"rating"`
}

type UpdateBook struct {
	Title      *string `json:"title"`
	Author     *string `json:"author"`
	Year       *int    `json:"year"`
	ISBN       *string `json:"ISBN"`
	OutOfStock *bool   `json:"outOfStock"`
	Rating     *int    `json:"rating"`
}
