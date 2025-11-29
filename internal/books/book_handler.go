package books

import (
	"context"
	"net/http"
	"strconv"
	"time"

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

func (h *BookHandler) ReadAllBooks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*3)
	defer cancel()

	pagePar := c.Query("page")
	title := c.Query("title")
	author := c.Query("author")
	category := c.Query("category")

	page, err := strconv.Atoi(pagePar)
	if err != nil {
		c.Error(err)
		return
	}

	filter := &Filters{
		Page:     int(page),
		Title:    title,
		Authors:  author,
		Category: category,
	}

	books, err := h.service.GetAll(ctx, filter)
	if err != nil {
		c.Error(middleware.InternalErr)
		c.Abort()
		return
	}

	response := make([]BookResponse, 0, len(books))
	for _, book := range books {
		response = append(response, ToResponse(&book))
	}

	c.JSON(http.StatusOK, response)
}

func (h *BookHandler) ReadBook(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*3)
	defer cancel()

	id, err := middleware.GetIdParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	book, err := h.service.GetById(ctx, id)
	if err != nil {
		c.Error(middleware.NotFound)
		return
	}

	c.JSON(http.StatusOK, ToResponse(book))
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	id, err := middleware.GetIdParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	var dtoReq BookRequest

	if err := c.ShouldBindJSON(&dtoReq); err != nil {
		c.Error(middleware.BadRequest)
		return
	}

	updateBook := &Books{
		ID:          id,
		Title:       dtoReq.Title,
		Description: dtoReq.Description,
		Content:     dtoReq.Content,
		AuthorID:    dtoReq.AuthorID,
	}

	if err := h.service.Update(c.Request.Context(), updateBook); err != nil {
		c.Error(middleware.InternalErr)
		return
	}

	c.JSON(http.StatusNoContent, "")
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	id, err := middleware.GetIdParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.Error(middleware.NotFound)
		return
	}

	c.JSON(http.StatusNoContent, "")
}

func (h *BookHandler) RelationBookCategory(c *gin.Context) {
	ctx := c.Request.Context()
	var bcDtoReq BookCategoryRequest

	if err := c.ShouldBindJSON(&bcDtoReq); err != nil {
		c.Error(middleware.BadRequest)
		return
	}

	if err := h.service.RelationBookCategory(ctx, bcDtoReq.BookID, bcDtoReq.CategoryID); err != nil {
		c.Error(middleware.InternalErr)
		return
	}

	c.JSON(http.StatusOK, "")
}
