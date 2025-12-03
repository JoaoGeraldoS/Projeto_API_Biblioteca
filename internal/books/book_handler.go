package books

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BookHandler struct {
	service BookServcie
	logApp  *zap.Logger
}

func NewBookHandler(svc BookServcie, logApp *zap.Logger) *BookHandler {
	return &BookHandler{service: svc, logApp: logApp}
}

// @Summary Cria um novo livro
// @Description Recebe um objeto JSON BookRequest e salva o livro no banco de dados.
// @Tags books
// @Accept  json
// @Produce json
// @Param   book body BookRequest true "Dados do Novo Livro a ser criado"
// @Success 201 {object} BookRequest "Livro criado com sucesso"
// @Failure 400 {object} middleware.APIError "Requisição Inválida (JSON malformado ou campo obrigatório ausente)"
// @Failure 500 {object} middleware.APIError "Erro interno do servidor"
// @Security ApiKeyAuth
// @Router /api/books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	h.logApp.Info("Rota de criar livros")

	ctx := c.Request.Context()

	var bDtoReq BookRequest

	if err := c.ShouldBindJSON(&bDtoReq); err != nil {
		h.logApp.Error("falha ao ler json", zap.Error(err))
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
		h.logApp.Error("falha ao criar livro", zap.Error(err))
		c.Error(middleware.InternalErr)
		return
	}

	c.JSON(http.StatusCreated, ToResponse(newBook))
}

// @Summary Listar livros
// @Description Retorna uma lista de livros filtrando por título, autor e categoria.
// @Tags books
// @Accept json
// @Produce json
// @Param page query int true "Página"
// @Param title query string false "Filtrar por título"
// @Param author query string false "Filtrar por autor"
// @Param category query string false "Filtrar por categoria"
// @Success 200 {array} BookResponse
// @Failure 400 {object} middleware.APIError "Parâmetros inválidos"
// @Failure 500 {object} middleware.APIError "Erro interno"
// @Router /public/api/books [get]
func (h *BookHandler) ReadAllBooks(c *gin.Context) {
	h.logApp.Info("Rota de ver todos os livros")

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*3)
	defer cancel()

	pagePar := c.Query("page")
	title := c.Query("title")
	author := c.Query("author")
	category := c.Query("category")

	page, err := strconv.Atoi(pagePar)
	if err != nil {
		h.logApp.Error("falha ao aonverter page", zap.Error(err))
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
		h.logApp.Error("falha ao obter livros", zap.Error(err))
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

// @Summary Obter livro
// @Description Retorna um livro específico.
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "ID do livro"
// @Success 200 {object} BookResponse
// @Failure 400 {object} middleware.APIError "ID inválido"
// @Failure 404 {object} middleware.APIError "Livro não encontrado"
// @Failure 500 {object}	middleware.APIError "Erro interno"
// @Router /public/api/books/{id} [get]
func (h *BookHandler) ReadBook(c *gin.Context) {
	h.logApp.Info("Rota de um livro")

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*3)
	defer cancel()

	id, err := middleware.GetIdParam(c)
	if err != nil {
		h.logApp.Error("falha ao verifica id", zap.Error(err))
		c.Error(err)
		return
	}

	book, err := h.service.GetById(ctx, id)
	if err != nil {
		h.logApp.Error("falha ao obter livro", zap.Error(err))
		c.Error(middleware.NotFound)
		return
	}

	c.JSON(http.StatusOK, ToResponse(book))
}

// @Summary Atualiza um livro
// @Description Recebe um objeto JSON BookRequest e atualiza o livro no banco de dados.
// @Tags books
// @Accept  json
// @Produce json
// @Param id path int true "Recebe o id do livro"
// @Param   book body BookRequest true "Dados do Novo Livro a ser atualizado"
// @Success 200 {object} BookRequest "Livro atualizado com sucesso"
// @Failure 400 {object} middleware.APIError "Requisição Inválida (JSON malformado ou campo obrigatório ausente)"
// @Failure 500 {object} middleware.APIError "Erro interno do servidor"
// @Security ApiKeyAuth
// @Router /api/books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	h.logApp.Info("Rota de atualizar um livro")

	id, err := middleware.GetIdParam(c)
	if err != nil {
		h.logApp.Error("falha ao verifica id", zap.Error(err))
		c.Error(err)
		return
	}

	var dtoReq BookRequest

	if err := c.ShouldBindJSON(&dtoReq); err != nil {
		h.logApp.Error("falha ao ler json", zap.Error(err))
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
		h.logApp.Error("falha ao atualizar livro", zap.Error(err))
		c.Error(middleware.InternalErr)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Exclui um livro pelo ID
// @Description Exclui um livro específico do banco de dados.
// @Tags books
// @Accept  json
// @Produce json
// @Security ApiKeyAuth
// @Param   id path int true "ID do Livro a ser excluído"
// @Success 204 "Nenhum Conteúdo"
// @Failure 400 {object} middleware.APIError "Requisição Inválida (ID com formato incorreto)"
// @Failure 404 {object} middleware.APIError "Livro não encontrado"
// @Router /api/books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	h.logApp.Info("Rota de apagar livro")

	id, err := middleware.GetIdParam(c)
	if err != nil {
		h.logApp.Error("falha ao verifica id", zap.Error(err))
		c.Error(err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.logApp.Error("falha ao apagar livro", zap.Error(err))
		c.Error(middleware.NotFound)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Associa um livro a uma categoria
// @Description Cria uma relação entre um livro existente e uma categoria existente.
// @Tags books
// @Accept  json
// @Produce json
// @Security ApiKeyAuth
// @Param   data body BookCategoryRequest true "IDs do Livro e da Categoria a serem relacionados"
// @Success 200 "OK, Relação criada com sucesso"
// @Failure 400 {object} middleware.APIError "Requisição Inválida (JSON malformado, campo ausente ou tipo errado)"
// @Failure 500 {object} middleware.APIError "Erro interno do servidor (Erro na criação da relação)"
// @Router /api/books/relation [post]
func (h *BookHandler) RelationBookCategory(c *gin.Context) {
	h.logApp.Info("Rota de relacionamento")
	ctx := c.Request.Context()
	var bcDtoReq BookCategoryRequest

	if err := c.ShouldBindJSON(&bcDtoReq); err != nil {
		h.logApp.Error("falha ao ler json", zap.Error(err))
		c.Error(middleware.BadRequest)
		return
	}

	if err := h.service.RelationBookCategory(ctx, bcDtoReq.BookID, bcDtoReq.CategoryID); err != nil {
		h.logApp.Error("falha ao fazer relacionamento", zap.Error(err))
		c.Error(middleware.InternalErr)
		return
	}

	c.Status(http.StatusOK)
}
