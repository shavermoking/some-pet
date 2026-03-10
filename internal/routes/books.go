package routes

import (
	"some-pet/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine, h *handlers.Books) {
	books := r.Group("/books")

	books.POST("/", h.Create)
	books.GET("/", h.GetAll)
	books.GET("/:id", h.GetByID)
	books.DELETE("/:id", h.Delete)
	books.PATCH("/:id", h.Update)
	books.POST("/:id/mark-out-of-stock", h.MarkOutOfStock)
	books.GET("/recommend", h.GetRecommend)
}
