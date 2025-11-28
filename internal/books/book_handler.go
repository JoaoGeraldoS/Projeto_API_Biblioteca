package books

import (
	"net/http"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/middleware"
	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	service BookServcie
}

func NewBookHandler(svc BookServcie) *BookHandler {
	return &BookHandler{service: svc}
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	ctx := c.Request.Context()

	var bDtoReq BookRequest

	if err := c.ShouldBindJSON(&bDtoReq); err != nil {
		c.Error(middleware.BadRequest)
		return
	}

	newBook := &Books{
		Title:       bDtoReq.Title,
		Description: bDtoReq.Description,
		Content:     bDtoReq.Content,
		AuthorID:    bDtoReq.AuthorID,
	}

	if err := h.service.Create(ctx, newBook); err != nil {
		c.Error(middleware.InternalErr)
		return
	}

	c.JSON(http.StatusCreated, newBook)
}

func (h *BookHandler) RelationBookCategory(c *gin.Context) {
	ctx := c.Request.Context()
	var bcDtoReq BookCategoryRequest

	if err := c.ShouldBindJSON(bcDtoReq); err != nil {
		c.Error(middleware.BadRequest)
		return
	}

	if err := h.service.RelationBookCategory(ctx, bcDtoReq.BookID, bcDtoReq.CategoryID); err != nil {
		c.Error(middleware.InternalErr)
		return
	}

	c.JSON(http.StatusOK, "")
}
