package service

import (
	"context"
	"some-pet/internal/models"
)

type BooksRepository interface {
	Create(ctx context.Context, book models.Book) (*models.Book, error)
	GetByID(ctx context.Context, id int) (models.Book, error)
	GetAll(ctx context.Context) ([]models.Book, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, input models.UpdateBook) error
	MarkOutOfStock(ctx context.Context, id int) error
	GetRecommend(ctx context.Context) ([]models.Book, error)
}

type Books struct {
	repo BooksRepository
}

func NewBooks(repo BooksRepository) *Books {
	return &Books{
		repo: repo,
	}
}

var ErrValidationFailed = "ошибка валидации"

func (b *Books) Create(ctx context.Context, book models.Book) (*models.Book, error) {
	return b.repo.Create(ctx, book)
}

func (b *Books) GetAll(ctx context.Context) ([]models.Book, error) {
	return b.repo.GetAll(ctx)
}

func (b *Books) GetByID(ctx context.Context, id int) (models.Book, error) {
	return b.repo.GetByID(ctx, id)
}

func (b *Books) Delete(ctx context.Context, id int) error {
	return b.repo.Delete(ctx, id)
}

func (b *Books) Update(ctx context.Context, id int, input models.UpdateBook) error {
	return b.repo.Update(ctx, id, input)
}

func (b *Books) MarkOutOfStock(ctx context.Context, id int) error {
	return b.repo.MarkOutOfStock(ctx, id)
}

func (b *Books) GetRecommend(ctx context.Context) ([]models.Book, error) {
	return b.repo.GetRecommend(ctx)
}
