package models

type Book struct {
	ID         int
	Title      string
	Author     string
	Year       int
	ISBN       string
	OutOfStock bool
	Rating     int
}
