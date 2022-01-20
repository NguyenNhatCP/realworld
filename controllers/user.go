package controllers

import (
	"crud-rest-api-golang/common"
	"crud-rest-api-golang/middlewares"
	"crud-rest-api-golang/models"
	"crud-rest-api-golang/serializers"
	"crud-rest-api-golang/validator"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func GetUsers(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var users []models.User
// 	database.DB.Find(&users)
// 	json.NewEncoder(w).Encode(users)
// }

// func GetUser(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	var user models.User
// 	database.DB.First(&user, params["id"])
// 	json.NewEncoder(w).Encode(user)
// }

// func CreateUser(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var user models.User
// 	json.NewDecoder(r.Body).Decode(&user)
// 	database.DB.Create(&user)
// 	json.NewEncoder(w).Encode(user)
// }
func UsersRegistration(c *gin.Context) {
	userModelValidator := validator.NewUserModelValidator()
	if err := userModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}
	userModel, _ := models.FindOneUser(&models.UserModel{Email: userModelValidator.UserModel.Email})
	if userModel.Email != "" {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("register", errors.New("Email: " + userModel.Email + " recent is existed."))) 
		return
	}
	if err := models.SaveOne(&userModelValidator.UserModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return 
	}
	c.Set("my_user_model", userModelValidator.UserModel)
	serializer := serializers.UserSerializer{C: c}
	c.JSON(http.StatusCreated, gin.H{"user": serializer.Response()})
}
// func CreateUser(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		var user models.UserModel
// 		json.NewDecoder(r.Body).Decode(&user)
// 		database.DB.Create(&user)
// 		json.NewEncoder(w).Encode(user)
// 	}

	func UserUpdate(c *gin.Context) {
		myUserModel := c.MustGet("my_user_model").(models.UserModel)
		userModelValidator := validator.NewUserModelValidatorFillWith(myUserModel)
		if err := userModelValidator.Bind(c); err != nil {
			c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
			return
		}

		userModelValidator.UserModel.ID = myUserModel.ID
		if err := models.UpdateUser((&userModelValidator.UserModel)); err != nil {
			c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
			return
		}
		middlewares.UpdateContextUserModel(c, myUserModel.ID)
		serializer := serializers.UserSerializer{C: c}
		c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
	}
//  func UpdateUser(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	var user models.User
// 	database.DB.First(&user, params["id"])
// 	json.NewDecoder(r.Body).Decode(&user)
// 	database.DB.Save(&user)
// 	json.NewEncoder(w).Encode(user)
// }

func UsersLogin(c *gin.Context) {
	loginValidator := validator.NewLoginValidator()
	if err := loginValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}
	userModel, err := models.FindOneUser(&models.UserModel{Email: loginValidator.UserModel.Email})

	if err != nil {
		c.JSON(http.StatusForbidden, common.NewError("login", errors.New("Not Registered email or invalid password")))
		return
	}

	if userModel.CheckPassword(loginValidator.User.Password) != nil {
		c.JSON(http.StatusForbidden, common.NewError("login", errors.New("Not Registered email or invalid password")))
		return
	}
	middlewares.UpdateContextUserModel(c, userModel.ID)
	serializer := serializers.UserSerializer{C: c}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}
// func DeleteUser(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	var user models.User
// 	database.DB.Delete(&user, params["id"])
// 	json.NewEncoder(w).Encode(user)
// }

func UserRetrieve(c *gin.Context) {
	serializer := serializers.UserSerializer{C: c}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}

//Profile
func ProfileRetrieve(c *gin.Context) {
	username := c.Param("username")
	userModel, err := models.FindOneUser(&models.UserModel{Username: username})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("profile", errors.New("Invalid username")))
		return
	}
	profileSerializer := serializers.ProfileSerializer{c, userModel}
	c.JSON(http.StatusOK, gin.H{"profile": profileSerializer.Response()})
}

func ProfileFollow(c *gin.Context) {
	username := c.Param("username")
	userModel, err := models.FindOneUser(&models.UserModel{Username: username})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("profile", errors.New("Invalid username")))
		return
	}
	myUserModel := c.MustGet("my_user_model").(models.UserModel)
	err = myUserModel.Following(userModel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := serializers.ProfileSerializer{c, userModel}
	c.JSON(http.StatusOK, gin.H{"profile": serializer.Response()})
}

func ProfileUnfollow(c *gin.Context) {
	username := c.Param("username")
	userModel, err := models.FindOneUser(&models.UserModel{Username: username})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("profile", errors.New("Invalid username")))
		return
	}
	myUserModel := c.MustGet("my_user_model").(models.UserModel)

	err = myUserModel.UnFollowing(userModel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := serializers.ProfileSerializer{c, userModel}
	c.JSON(http.StatusOK, gin.H{"profile": serializer.Response()})
}