package categories

import (
	"net/http"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/middleware"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	svc CategoryService
}

func NewCategoryHandler(svc CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	ctx := c.Request.Context()

	var dto CategoryRequest

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(middleware.BadRequest)
		c.Abort()
		return
	}

	category := &Category{
		Name: dto.Name,
	}

	if err := h.svc.Create(ctx, category); err != nil {
		c.Error(middleware.InternalErr)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, ToResponse(category))
}

func (h *CategoryHandler) ReadCategories(c *gin.Context) {
	allCategories, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		c.Error(middleware.InternalErr)
		c.Abort()
		return
	}

	response := make([]CategoryResponse, 0, len(allCategories))
	for _, cat := range allCategories {
		response = append(response, ToResponse(&cat))
	}

	c.JSON(http.StatusOK, response)
}

func (h *CategoryHandler) ReadCategory(c *gin.Context) {
	id, err := middleware.GetIdParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	category, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		c.Error(middleware.NotFound)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, ToResponse(category))
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := middleware.GetIdParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	var dto CategoryRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(middleware.BadRequest)
		c.Abort()
		return
	}

	updateCategory := &Category{
		ID:   id,
		Name: dto.Name,
	}

	err = h.svc.Update(c.Request.Context(), updateCategory)
	if err != nil {
		c.Error(middleware.InternalErr)
		return
	}

	c.JSON(http.StatusNoContent, "")
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := middleware.GetIdParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.Error(middleware.NotFound)
	}

	c.JSON(http.StatusNoContent, "")
}
