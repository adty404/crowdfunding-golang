package handler

import (
	"crowdfunding-golang/helper"
	"crowdfunding-golang/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	service user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{service: userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	newUser, err := h.service.RegisterUser(input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", newUser)

	c.JSON(http.StatusOK, response)
}