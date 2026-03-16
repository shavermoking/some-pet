package handlers

import (
	"errors"
	"net/http"
	"some-pet/internal/models"
	"some-pet/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Books struct {
	service *service.Books
}

func NewBooks(service *service.Books) *Books {
	return &Books{service: service}
}

func (h *Books) Create(c *gin.Context) {
	var input models.CreateBook

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректный формат данных: " + err.Error(),
		})
		return
	}

	book := models.Book{
		Title:      input.Title,
		Author:     input.Author,
		Year:       input.Year,
		ISBN:       input.ISBN,
		Rating:     input.Rating,
		OutOfStock: false,
	}

	created, err := h.service.Create(c.Request.Context(), book)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "книга с таким ISBN уже существует или данные некорректны",
		})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *Books) GetAll(c *gin.Context) {
	books, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (h *Books) GetByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	book, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *Books) Delete(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.Status(http.StatusNoContent)
}

func (h *Books) Update(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	var input models.UpdateBook

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.service.Update(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Books) MarkOutOfStock(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	err = h.service.MarkOutOfStock(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "book not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Books) GetRecommend(c *gin.Context) {
	books, err := h.service.GetRecommend(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, books)
}
