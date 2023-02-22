package main

import (
	"crowdfunding-golang/auth"
	"crowdfunding-golang/handler"
	"crowdfunding-golang/helper"
	"crowdfunding-golang/user"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:root@tcp(localhost:3306)/crowdfunding_golang?charset=utf8&parseTime=True&loc=Local"
	// dsn := "root:@tcp(localhost:3306)/crowdfunding_golang?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewService()
	userHandler := handler.NewUserHandler(userService, authService)

	router := gin.Default()
	api := router.Group("/api/v1") // API Versioning

	// Route
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)

	router.Run()
}

func authMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if !strings.Contains(authHeader, "Bearer") {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	tokenString := ""
	arrayToken := strings.Split(authHeader, " ")
	if len(arrayToken) == 2 {
		tokenString = arrayToken[1]
	}


}
