package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/middleware"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc UserService
}

func NewUsersHandler(svc UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// @Summary Cria um novo usuario
// @Description Recebe um objeto JSON UserRequest e salva o usuario no banco de dados.
// @Tags users
// @Accept  json
// @Produce json
// @Param   user body UserRequest true "Dados do Novo usuario a ser criado"
// @Success 201 {object} UserRequest "Usuario criado com sucesso"
// @Failure 400 {object} middleware.APIError "Requisição Inválida (JSON malformado ou campo obrigatório ausente)"
// @Failure 500 {object} middleware.APIError "Erro interno do servidor"
// @Security ApiKeyAuth
// @Router /public/api/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {

	var dtoReq UserRequest

	if err := c.ShouldBindJSON(&dtoReq); err != nil {
		fmt.Println(err)
		c.Error(middleware.BadRequest)
		return
	}

	newUser := &Users{
		Name:     dtoReq.Name,
		Email:    dtoReq.Email,
		Username: dtoReq.Username,
		Password: dtoReq.Password,
		Role:     dtoReq.Role,
	}

	if err := h.svc.Create(c.Request.Context(), newUser); err != nil {
		c.Error(middleware.InternalErr)
		return
	}

	c.JSON(http.StatusCreated, ToResponse(newUser))

}

// @Summary Listar usuarios
// @Description Retorna uma lista de usuarios
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} UserResponse
// @Failure 500 {object} middleware.APIError "Erro interno"
// @Router /public/api/users [get]
func (h *UserHandler) ReadAllUsers(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	getUsers, err := h.svc.GetAll(ctx)
	if err != nil {
		c.Error(middleware.InternalErr)
		return
	}

	response := make([]UserResponse, 0, len(getUsers))
	for _, user := range getUsers {
		response = append(response, ToResponse(&user))
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Obter usuario
// @Description Retorna um usuario
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "Recebe o id do usuario"
// @Success 200 {object} UserResponse
// @Failure 500 {object} middleware.APIError "Erro interno"
// @Router /public/api/users/{id} [get]
func (h *UserHandler) ReadUser(c *gin.Context) {
	id, err := middleware.GetIdParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	result, err := h.svc.GetById(ctx, id)
	if err != nil {
		fmt.Println(err)
		c.Error(middleware.NotFound)
		return
	}

	c.JSON(http.StatusOK, ToResponse(result))
}

// @Summary Atualiza um usuario
// @Description Recebe um objeto JSON UserRequest e atualiza o usuario no banco de dados.
// @Tags users
// @Accept  json
// @Produce json
// @Param id path int true "Recebe o id da usuario"
// @Param   user body UserRequest true "Dados do novo usuario a ser atualizado"
// @Success 204  "Usuario atualizado com sucesso"
// @Failure 400 {object} middleware.APIError "Requisição Inválida (JSON malformado ou campo obrigatório ausente)"
// @Failure 500 {object} middleware.APIError "Erro interno do servidor"
// @Security ApiKeyAuth
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := middleware.GetIdParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	var dtoReq UserRequest

	if err := c.ShouldBindJSON(&dtoReq); err != nil {
		c.Error(middleware.BadRequest)
		return
	}

	updateUser := &Users{
		ID:   id,
		Name: dtoReq.Name,
		Bio:  dtoReq.Bio,
	}

	if err := h.svc.Update(c.Request.Context(), updateUser); err != nil {
		c.Error(middleware.InternalErr)
		return
	}

	c.JSON(http.StatusNoContent, "")
}

// @Summary Exclui um usuario pelo ID
// @Description Exclui um usuario específico do banco de dados.
// @Tags users
// @Accept  json
// @Produce json
// @Security ApiKeyAuth
// @Param   id path int true "ID do usuario a ser excluído"
// @Success 204 "Nenhum Conteúdo"
// @Failure 400 {object} middleware.APIError "Requisição Inválida (ID com formato incorreto)"
// @Failure 404 {object} middleware.APIError "Livro não encontrado"
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
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

// @Summary Faz o login do usuario
// @Description Recebe um objeto JSON LoginRequest.
// @Tags users
// @Accept  json
// @Produce json
// @Param   user body LoginRequest true "Dados para realizar o login"
// @Success 200 "token do usuario"
// @Failure 400 {object} middleware.APIError "Requisição Inválida (JSON malformado ou campo obrigatório ausente)"
// @Failure 500 {object} middleware.APIError "Erro interno do servidor"
// @Router /public/api/users/login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	var dtoLogin LoginRequest

	if err := c.ShouldBindJSON(&dtoLogin); err != nil {
		c.Error(middleware.BadRequest)
		return
	}

	user, err := h.svc.Login(ctx, dtoLogin.Email, dtoLogin.Password)
	if err != nil {
		fmt.Println(err)
		c.Error(middleware.InternalErr)
		return
	}

	tokenString, err := middleware.GenerateToken(user.Email, string(user.Role))
	if err != nil {
		c.Error(middleware.InternalErr)
		return
	}

	c.JSON(http.StatusOK, tokenString)
}
