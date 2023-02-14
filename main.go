package main

import (
	"crowdfunding-golang/user"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:root@tcp(localhost:3306)/crowdfunding_golang?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	userRepository := user.NewRepository(db)
	user := user.User{
		Name:       "John",
		Occupation: "Programmer",
	}

	_, err = userRepository.Save(user)
}
