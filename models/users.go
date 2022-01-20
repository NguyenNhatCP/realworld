package models

import (
	"crud-rest-api-golang/database"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var DB *gorm.DB

type UserModel struct {
	ID           uint    `gorm:"primary_key"`
	Username     string  `gorm:"column:username"`
	Email        string  `gorm:"column:email;unique"`
	Bio          string  `gorm:"column:bio;size:1024"`
	Image        *string `gorm:"column:image"`
	PasswordHash string  `gorm:"column:password;not null"`
	Followers  [] FollowModel  `gorm:"foreignkey:FollowingID"`
	Followings [] FollowModel  `gorm:"foreignkey:FollowedByID"`
}

// A hack way to save ManyToMany relationship,
type FollowModel struct {
	gorm.Model
	Following    UserModel
	FollowingID  uint
	FollowedBy   UserModel
	FollowedByID uint
}

func AutoMigrate() {
	db := database.GetDB()

	db.AutoMigrate(&UserModel{})
	db.AutoMigrate(&FollowModel{})
}

func (u *UserModel) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty!")
	}
	bytePassword := []byte(password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.PasswordHash = string(passwordHash)
	return nil
}

func (u *UserModel) CheckPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.PasswordHash)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

func SaveOne(data interface{}) error {
	err := database.GetDB().Save(data).Error
	return err
}
func UpdateUser(data interface{}) error {
	db := database.GetDB()
	err := db.Save(data).Error
	return err
}
func FindOneUser(condition interface{}) (UserModel, error) {
	var model UserModel
	err := database.GetDB().Where(condition).First(&model).Error
	return model, err
}

func (u UserModel) IsFollowing(v UserModel) bool {
	db := database.GetDB()
	var follow FollowModel
	db.Where(FollowModel{
		FollowingID:  v.ID,
		FollowedByID: u.ID,
	}).First(&follow)
	return follow.ID != 0
}

func (u UserModel) Following(v UserModel) error {
	db := database.GetDB()
	var follow FollowModel
	err := db.FirstOrCreate(&follow, &FollowModel{
		FollowingID:  v.ID,
		FollowedByID: u.ID,
	}).Error
	return err
}

func (u UserModel) UnFollowing(v UserModel) error {
	db := database.GetDB()
	err := db.Where(&FollowModel{
		FollowingID:  v.ID,
		FollowedByID: u.ID,
	}).Delete(&FollowModel{}).Error
	return err
}

// You could get a following list of userModel
// 	followings := userModel.GetFollowings()
func (u UserModel) GetFollowings() []UserModel {
	db := database.GetDB()
	tx := db.Begin()
	var follows []FollowModel
	var followings []UserModel
	tx.Debug().Where(FollowModel{
		FollowedByID: u.ID,
	}).Preload("Following").Preload("FollowedBy").Find(&follows)
	for _, follow := range follows {
		var userModel UserModel
		tx.Debug().Where("id = (?)", follow.Following.ID).Find(&userModel)
		followings = append(followings, userModel)
	}
	tx.Commit()
	return followings
}
