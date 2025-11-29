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

func (h *UserHandler) CreateUser(c *gin.Context) {

	var dtoReq UserRequest

	if err := c.ShouldBindJSON(&dtoReq); err != nil {
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
		c.Error(middleware.NotFound)
		return
	}

	c.JSON(http.StatusOK, ToResponse(result))
}

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
