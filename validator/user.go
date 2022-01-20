package validator

import (
	"crud-rest-api-golang/common"
	"crud-rest-api-golang/models"

	"github.com/gin-gonic/gin"
)

type UserModelValidator struct {
	User struct {
		Username string `json:"username" binding:"required,min=4"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8,max=255"`
		Bio      string `json:"bio" binding:"max=1024"`
		Image    string `json:"image" binding:"omitempty,url"`
	} `json:"user"`
	UserModel models.UserModel `json:"-"`
}
// Init default
func NewUserModelValidator() UserModelValidator {
	userModelValidator := UserModelValidator{}
	return userModelValidator
}
func NewUserModelValidatorFillWith(userModel models.UserModel) UserModelValidator {
	userModelValidator := NewUserModelValidator()
	userModelValidator.User.Username = userModel.Username
	userModelValidator.User.Email = userModel.Email
	userModelValidator.User.Bio = userModel.Bio
	userModelValidator.User.Password = common.NBRandomPassword

	if userModel.Image != nil {
		userModelValidator.User.Image = *userModel.Image
	}
	return userModelValidator
}

type LoginValidator struct {
	User struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8,max=255"`
	} `json:"user"`
	UserModel models.UserModel `json:"-"`
}

func (login *LoginValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, login)
	if err != nil {
		return err
	}

	login.UserModel.Email = login.User.Email
	return nil
}

// Init default
func NewLoginValidator() LoginValidator {
	loginValidator := LoginValidator{}
	return loginValidator
}


func (model *UserModelValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, model)
	if err != nil {
		return err
	}
	model.UserModel.Username = model.User.Username
	model.UserModel.Email = model.User.Email
	model.UserModel.Bio = model.User.Bio

	if model.User.Password != common.NBRandomPassword {
		model.UserModel.SetPassword(model.User.Password)
	}
	if model.User.Image != "" {
		model.UserModel.Image = &model.User.Image
	}
	return nil
}
