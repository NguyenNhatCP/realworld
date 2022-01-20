package main

import (
	"crud-rest-api-golang/database"
	"crud-rest-api-golang/middlewares"
	"crud-rest-api-golang/models"
	"crud-rest-api-golang/router"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
func Migrate(db *gorm.DB) {
	models.AutoMigrate()
	db.AutoMigrate(&models.ArticleModel{})
	db.AutoMigrate(&models.TagModel{})
	db.AutoMigrate(&models.FavoriteModel{})
	db.AutoMigrate(&models.ArticleUserModel{})
	db.AutoMigrate(&models.CommentModel{})
}
func main() {

	// init db
	db := database.InitialDB()

	Migrate(db)

	// Close
	defer database.CloseDB()

	// init port getway
	r := gin.Default()

	v1 := r.Group("/api")
	v1.Use(middlewares.AuthMiddleware(false))
	router.UsersRegister(v1.Group("/users"))
	router.ArticlesAnonymousRegister(v1.Group("/articles"))

	v1.Use(middlewares.AuthMiddleware(true))
	router.UserRegister(v1.Group("/user"))
	router.ProfileRegister(v1.Group("/profiles"))

	// acticle
	router.ArticlesRegister(v1.Group("/articles"))

	// tag
	router.TagsAnonymousRegister(v1.Group("/tags"))
	r.Run()
	// router.InitializeRouter()
}
