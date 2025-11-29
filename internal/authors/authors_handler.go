package authors

import (
	"net/http"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/middleware"
	"github.com/gin-gonic/gin"
)

type AuthorHandler struct {
	svc AuthorsService
}

func NewAuthorsHandler(svc AuthorsService) *AuthorHandler {
	return &AuthorHandler{svc: svc}
}

func (h *AuthorHandler) CreateAuthor(c *gin.Context) {
	ctx := c.Request.Context()

	var dto AuthorRequest

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(middleware.BadRequest)
		c.Abort()
		return
	}

	author := &Authors{
		Name:        dto.Name,
		Description: dto.Description,
	}

	if err := h.svc.Create(ctx, author); err != nil {
		c.Error(middleware.InternalErr)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, ToResponse(author))
}

func (h *AuthorHandler) ReadAuthors(c *gin.Context) {
	authors, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		c.Error(middleware.InternalErr)
		return
	}

	response := make([]AuthorResponse, 0, len(authors))
	for _, a := range authors {
		response = append(response, ToResponse(&a))
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthorHandler) ReadAuthor(c *gin.Context) {
	id, err := middleware.GetIdParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	author, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.Error(middleware.NotFound)
		return
	}

	c.JSON(http.StatusOK, ToResponse(author))
}

func (h *AuthorHandler) UpdateAuthor(c *gin.Context) {
	id, err := middleware.GetIdParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	var dto AuthorRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(middleware.BadRequest)
		return
	}

	updateAuthor := &Authors{
		ID:          id,
		Name:        dto.Name,
		Description: dto.Description,
	}

	err = h.svc.Update(c.Request.Context(), updateAuthor)
	if err != nil {
		c.Error(middleware.InternalErr)
		return
	}

	c.JSON(http.StatusNoContent, "")
}

func (h *AuthorHandler) DeleteAuthor(c *gin.Context) {
	id, err := middleware.GetIdParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.Error(middleware.NotFound)
		return
	}

	c.JSON(http.StatusNoContent, "")

}
