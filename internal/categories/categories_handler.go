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

// @Summary Cria um nova categoria
// @Description Recebe um objeto JSON CategoryRequest e salva a categoria no banco de dados.
// @Tags categories
// @Accept  json
// @Produce json
// @Param   category body CategoryRequest true "Dados da Nova categoria a ser criado"
// @Success 201 {object} CategoryRequest "Categoria criada com sucesso"
// @Failure 400 {object} middleware.APIError "Requisição Inválida (JSON malformado ou campo obrigatório ausente)"
// @Failure 500 {object} middleware.APIError "Erro interno do servidor"
// @Security ApiKeyAuth
// @Router /api/categories [post]
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

// @Summary Listar categorias
// @Description Retorna uma lista de categorias
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} CategoryResponse
// @Failure 500 {object} middleware.APIError "Erro interno"
// @Router /public/api/categories [get]
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

// @Summary Obtem categoria
// @Description Retorna uma categorias
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Recebe o id da categoria"
// @Success 200 {object} CategoryResponse
// @Failure 500 {object} middleware.APIError "Erro interno"
// @Router /public/api/categories/{id} [get]
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

// @Summary Atualiza uma categoria
// @Description Recebe um objeto JSON CategoryRequest e atualiza a categoria no banco de dados.
// @Tags categories
// @Accept  json
// @Produce json
// @Param id path int true "Recebe o id da categoria"
// @Param   category body CategoryRequest true "Dados da nova categoria a ser atualizada"
// @Success 204  "Categoria atualizada com sucesso"
// @Failure 400 {object} middleware.APIError "Requisição Inválida (JSON malformado ou campo obrigatório ausente)"
// @Failure 500 {object} middleware.APIError "Erro interno do servidor"
// @Security ApiKeyAuth
// @Router /api/categories/{id} [put]
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

// @Summary Exclui uma categoria pelo ID
// @Description Exclui uma categoria específica do banco de dados.
// @Tags categories
// @Accept  json
// @Produce json
// @Security ApiKeyAuth
// @Param   id path int true "ID da categoria a ser excluída"
// @Success 204 "Nenhum Conteúdo"
// @Failure 400 {object} middleware.APIError "Requisição Inválida (ID com formato incorreto)"
// @Failure 404 {object} middleware.APIError "Livro não encontrado"
// @Router /api/categories/{id} [delete]
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
