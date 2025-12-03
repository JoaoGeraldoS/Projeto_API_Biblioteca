package authors

import (
	"net/http"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthorHandler struct {
	svc    AuthorsService
	logApp *zap.Logger
}

func NewAuthorsHandler(svc AuthorsService, log *zap.Logger) *AuthorHandler {
	return &AuthorHandler{svc: svc, logApp: log}
}

// @Summary Cria um novo autor
// @Description Recebe um objeto JSON AuthorRequest e salva o autor no banco de dados.
// @Tags authors
// @Accept  json
// @Produce json
// @Param   authors body AuthorRequest true "Dados do novo autor a ser criado"
// @Success 201 {object} AuthorRequest "Autor criada com sucesso"
// @Failure 400 {object} middleware.APIError "Requisição Inválida (JSON malformado ou campo obrigatório ausente)"
// @Failure 500 {object} middleware.APIError "Erro interno do servidor"
// @Security ApiKeyAuth
// @Router /api/authors [post]
func (h *AuthorHandler) CreateAuthor(c *gin.Context) {
	h.logApp.Info("Rode de criar autor")

	ctx := c.Request.Context()

	var dto AuthorRequest

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.logApp.Error("falha ao ler json", zap.Error(err))
		c.Error(middleware.BadRequest)
		c.Abort()
		return
	}

	author := &Authors{
		Name:        dto.Name,
		Description: dto.Description,
	}

	if err := h.svc.Create(ctx, author); err != nil {
		h.logApp.Error("falha ao criar autor", zap.Error(err))
		c.Error(middleware.InternalErr)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, ToResponse(author))
}

// @Summary Listar autores
// @Description Retorna uma lista de autores
// @Tags authors
// @Accept json
// @Produce json
// @Success 200 {array} AuthorResponse
// @Failure 500 {object} middleware.APIError "Erro interno"
// @Router /public/api/authors [get]
func (h *AuthorHandler) ReadAuthors(c *gin.Context) {
	h.logApp.Info("Rode de obter autores")

	authors, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		h.logApp.Error("falha ao obter autores", zap.Error(err))
		c.Error(middleware.InternalErr)
		return
	}

	response := make([]AuthorResponse, 0, len(authors))
	for _, a := range authors {
		response = append(response, ToResponse(&a))
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Obter autor
// @Description Retorna um autor
// @Tags authors
// @Accept json
// @Produce json
// @Param id path int true "Recebe o id do autor"
// @Success 200 {object} AuthorResponse
// @Failure 500 {object} middleware.APIError "Erro interno"
// @Router /public/api/authors/{id} [get]
func (h *AuthorHandler) ReadAuthor(c *gin.Context) {
	h.logApp.Info("Roda de obter autor")

	id, err := middleware.GetIdParam(c)
	if err != nil {
		h.logApp.Error("falha ao verificar id", zap.Error(err))
		c.Error(err)
		return
	}

	author, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logApp.Error("falha ao obter autor", zap.Error(err))
		c.Error(middleware.NotFound)
		return
	}

	c.JSON(http.StatusOK, ToResponse(author))
}

// @Summary Atualiza um autor
// @Description Recebe um objeto JSON AuthorRequest e atualiza o autor no banco de dados.
// @Tags authors
// @Accept  json
// @Produce json
// @Param id path int true "Recebe o id do autor"
// @Param   authors body AuthorRequest true "Dados do Novo autor a ser atualizado"
// @Success 204  "Autor atualizado com sucesso"
// @Failure 400 {object} middleware.APIError "Requisição Inválida (JSON malformado ou campo obrigatório ausente)"
// @Failure 500 {object} middleware.APIError "Erro interno do servidor"
// @Security ApiKeyAuth
// @Router /api/authors/{id} [put]
func (h *AuthorHandler) UpdateAuthor(c *gin.Context) {
	h.logApp.Info("Rota de atualizar autor")

	id, err := middleware.GetIdParam(c)
	if err != nil {
		h.logApp.Error("falha ao verificar id", zap.Error(err))
		c.Error(err)
		return
	}

	var dto AuthorRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		h.logApp.Error("falha ao ler json", zap.Error(err))
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
		h.logApp.Error("erro ao atualizar autor", zap.Error(err))
		c.Error(middleware.InternalErr)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Exclui um autor pelo ID
// @Description Exclui um autor específico do banco de dados.
// @Tags authors
// @Accept  json
// @Produce json
// @Security ApiKeyAuth
// @Param   id path int true "ID do autor a ser excluído"
// @Success 204 "Nenhum Conteúdo"
// @Failure 400 {object} middleware.APIError "Requisição Inválida (ID com formato incorreto)"
// @Failure 404 {object} middleware.APIError "Livro não encontrado"
// @Router /api/authors/{id} [delete]
func (h *AuthorHandler) DeleteAuthor(c *gin.Context) {
	h.logApp.Info("Rota de apagar autor")

	id, err := middleware.GetIdParam(c)
	if err != nil {
		h.logApp.Error("falha ao verificar id", zap.Error(err))
		c.Error(err)
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		h.logApp.Error("falha ao apagar autor", zap.Error(err))
		c.Error(middleware.NotFound)
		return
	}

	c.Status(http.StatusNoContent)
}
