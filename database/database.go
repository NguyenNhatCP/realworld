package database

import (
	"crud-rest-api-golang/env"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

func InitialDB() *gorm.DB  { 
	db, err := gorm.Open(mysql.Open(env.ConnectionString), &gorm.Config{})

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(env.ConnectionString)
		panic("Can't connect to DB!")
	}
	//db.LogMode(true)
	DB = db
	return DB

	// // init user
	// DB.AutoMigrate(&models.User{})
	// DB.AutoMigrate(&models.UserModel{})
	// DB.AutoMigrate(&models.FollowModel{})

	// // init article
	// DB.AutoMigrate(&models.ArticleModel{})
	// DB.AutoMigrate(&models.ArticleUserModel{})
	// DB.AutoMigrate(&models.FavoriteModel{})
	// DB.AutoMigrate(&models.TagModel{})
	// DB.AutoMigrate(&models.CommentModel{})
}
func GetDB() *gorm.DB {
	return DB
}

func CloseDB(){
	sqlDB, err := DB.DB();
	if err != nil {
		fmt.Println(err.Error())
		panic("Can't close to DB!")
	}
	sqlDB.Close()
}