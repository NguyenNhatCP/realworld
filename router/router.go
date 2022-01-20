package router

import (
	"crud-rest-api-golang/controllers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// router user
func UsersRegister(router *gin.RouterGroup) {
	router.POST("/", controllers.UsersRegistration)
	router.POST("/login", controllers.UsersLogin)
}
func UserRegister(router *gin.RouterGroup) {
	router.GET("/", controllers.UserRetrieve)
	router.PUT("/", controllers.UserUpdate)
}

func ProfileRegister(router *gin.RouterGroup) {
	router.GET("/:username", controllers.ProfileRetrieve)
	router.POST("/:username/follow", controllers.ProfileFollow)
	router.DELETE("/:username/follow", controllers.ProfileUnfollow)
}


// acticle
func ArticlesRegister(router *gin.RouterGroup) {
	router.POST("/", controllers.ArticleCreate)
	router.PUT("/:slug", controllers.ArticleUpdate)
	router.DELETE("/:slug", controllers.ArticleDelete)
	router.POST("/:slug/favorite", controllers.ArticleFavorite)
	router.DELETE("/:slug/favorite", controllers.ArticleUnfavorite)
	router.POST("/:slug/comments", controllers.ArticleCommentCreate)
	router.DELETE("/:slug/comments/:id", controllers.ArticleCommentDelete)
}

func ArticlesAnonymousRegister(router *gin.RouterGroup) {
	router.GET("/", controllers.ArticleList)
	router.GET("/:slug", controllers.ArticleRetrieve)
	router.GET("/:slug/comments", controllers.ArticleCommentList)
}

// tag
func TagsAnonymousRegister(router *gin.RouterGroup) {
	router.GET("/", controllers.TagList)
}
func InitializeRouter() {
	r := mux.NewRouter()
	
	// r.HandleFunc("/api/users", controllers.GetUsers).Methods("GET")
	// r.HandleFunc("/api/users/{id}", controllers.GetUser).Methods("GET")
	// r.HandleFunc("/api/users", controllers.CreateUser).Methods("POST")
	// r.HandleFunc("/api/users", controllers.UsersRegistration).Methods("POST")
	// r.HandleFunc("/api/login", controllers.CreateUser).Methods("POST")
	// r.HandleFunc("/api/users/{id}", controllers.UpdateUser).Methods("PUT")
	// r.HandleFunc("/api/users/{id}", controllers.DeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9000",
		handlers.CORS(handlers.AllowedHeaders([]string{
			"X-Requested-With",
			"Content-Type",
			"Authorization"}),
			handlers.AllowedMethods([]string{
				"GET",
				"POST",
				"PUT",
				"DELETE",
				"HEAD",
				"OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(r)))

}