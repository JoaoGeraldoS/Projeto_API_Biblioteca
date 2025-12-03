package categories

import (
	"net/http"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	svc    CategoryService
	logApp *zap.Logger
}

func NewCategoryHandler(svc CategoryService, log *zap.Logger) *CategoryHandler {
	return &CategoryHandler{svc: svc, logApp: log}
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
	h.logApp.Info("Rota de criar categoria")

	ctx := c.Request.Context()

	var dto CategoryRequest

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.logApp.Error("falha ao ler json", zap.Error(err))
		_ = c.Error(middleware.BadRequest)
		c.Abort()
		return
	}

	category := &Category{
		Name: dto.Name,
	}

	if err := h.svc.Create(ctx, category); err != nil {
		h.logApp.Error("falha ao criar categoria", zap.Error(err))
		_ = c.Error(middleware.InternalErr)
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
	h.logApp.Info("Rota de obter categorias")

	allCategories, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		h.logApp.Error("falha ao obter categorias", zap.Error(err))
		_ = c.Error(middleware.InternalErr)
		c.Abort()
		return
	}

	response := make([]CategoryResponse, 0, len(allCategories))
	for _, cat := range allCategories {
		response = append(response, ToResponse(&cat))
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Obter categoria
// @Description Retorna uma categorias
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Recebe o id da categoria"
// @Success 200 {object} CategoryResponse
// @Failure 500 {object} middleware.APIError "Erro interno"
// @Router /public/api/categories/{id} [get]
func (h *CategoryHandler) ReadCategory(c *gin.Context) {
	h.logApp.Info("Rota de obter categoria")

	id, err := middleware.GetIdParam(c)
	if err != nil {
		h.logApp.Error("falha ao verificar id", zap.Error(err))
		_ = c.Error(err)
		return
	}

	category, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		h.logApp.Error("falha ao obter categoria", zap.Error(err))
		_ = c.Error(middleware.NotFound)
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
	h.logApp.Info("Rota de atualizar categoria")

	id, err := middleware.GetIdParam(c)
	if err != nil {
		h.logApp.Error("falha ao verificar id", zap.Error(err))
		_ = c.Error(err)
		return
	}

	var dto CategoryRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		h.logApp.Error("falha ao ler json", zap.Error(err))
		_ = c.Error(middleware.BadRequest)
		c.Abort()
		return
	}

	updateCategory := &Category{
		ID:   id,
		Name: dto.Name,
	}

	err = h.svc.Update(c.Request.Context(), updateCategory)
	if err != nil {
		h.logApp.Error("falha ao atualizar categoria", zap.Error(err))
		_ = c.Error(middleware.InternalErr)
		return
	}

	c.Status(http.StatusNoContent)
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
	h.logApp.Info("Rota de apagar categoria")

	id, err := middleware.GetIdParam(c)
	if err != nil {
		h.logApp.Error("falha ao verificar id", zap.Error(err))
		_ = c.Error(err)
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		h.logApp.Error("falha ao apagar categoria", zap.Error(err))
		_ = c.Error(middleware.NotFound)
		return
	}

	c.Status(http.StatusNoContent)
}
