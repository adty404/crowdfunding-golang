package main

import (
	"crowdfunding-golang/handler"
	"crowdfunding-golang/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(localhost:3306)/crowdfunding_golang?charset=utf8&parseTime=True&loc=Local"
	// dsn := "root:root@tcp(localhost:3306)/crowdfunding_golang?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	userByEmail, err := userRepository.FindByEmail("adty404@gmail.com")

	if err != nil {
		log.Fatal(err)
	}

	if (user.User{}) == userByEmail {
		log.Println("User not found")
	}

	log.Println(userByEmail)

	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	api := router.Group("/api/v1") // API Versioning

	// Route
	api.POST("/users", userHandler.RegisterUser)

	router.Run()
}
