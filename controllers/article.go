package controllers

import (
	"crud-rest-api-golang/common"
	"crud-rest-api-golang/models"
	"crud-rest-api-golang/serializers"
	"crud-rest-api-golang/validator"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)
func ArticleCreate(c *gin.Context) {
	articleModelValidator := validator.NewArticleModelValidator()
	if err := articleModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}
	fmt.Println(articleModelValidator.ArticleModel.Author.UserModel)

	if err := models.SaveOne(&articleModelValidator.ArticleModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := serializers.ArticleSerializer{c, articleModelValidator.ArticleModel}
	c.JSON(http.StatusCreated, gin.H{"article": serializer.Response()})
}
func ArticleUpdate(c *gin.Context) {
	slug := c.Param("slug")
	articleModel, err := models.FindOneArticle(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("articles", errors.New("Invalid slug")))
		return
	}
	articleModelValidator := validator.NewArticleModelValidatorFillWith(articleModel)
	if err := articleModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	articleModelValidator.ArticleModel.ID = articleModel.ID
	if err := articleModel.Update(&articleModelValidator.ArticleModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := serializers.ArticleSerializer{c, articleModelValidator.ArticleModel}
	c.JSON(http.StatusOK, gin.H{"article": serializer.Response()})
}

func ArticleDelete(c *gin.Context) {
	slug := c.Param("slug")
	model, err := models.FindOneArticle(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("articles", errors.New("Invalid slug")))
		return
	}
	err = models.DeleteArticleModel(model)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"article": "Delete article success"})
}


func ArticleFavorite(c *gin.Context) {
	slug := c.Param("slug")
	articleModel, err := models.FindOneArticle(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("articles", errors.New("Invalid slug")))
		return
	}
	myUserModel := c.MustGet("my_user_model").(models.UserModel)
	err = articleModel.FavoriteBy(models.GetArticleUserModel(myUserModel))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := serializers.ArticleSerializer{c, articleModel}
	c.JSON(http.StatusOK, gin.H{"article": serializer.Response()})
}

func ArticleUnfavorite(c *gin.Context) {
	slug := c.Param("slug")
	articleModel, err := models.FindOneArticle(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("articles", errors.New("Invalid slug")))
		return
	}
	myUserModel := c.MustGet("my_user_model").(models.UserModel)
	err = articleModel.UnFavoriteBy(models.GetArticleUserModel(myUserModel))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := serializers.ArticleSerializer{c, articleModel}
	c.JSON(http.StatusOK, gin.H{"article": serializer.Response()})
}

func ArticleCommentCreate(c *gin.Context) {
	slug := c.Param("slug")
	articleModel, err := models.FindOneArticle(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("comment", errors.New("Invalid slug")))
		return
	}
	commentModelValidator := validator.NewCommentModelValidator()
	if err := commentModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}
	commentModelValidator.CommentModel.Article = articleModel

	if err := models.SaveOne(&commentModelValidator.CommentModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := serializers.CommentSerializer{c, commentModelValidator.CommentModel}
	c.JSON(http.StatusCreated, gin.H{"comment": serializer.Response()})
}

func ArticleCommentDelete(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	id := uint(id64)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("comment", errors.New("Invalid id")))
		return
	}
	err = models.DeleteCommentModel([]uint{id})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("comment", errors.New("Invalid id")))
		return
	}
	c.JSON(http.StatusOK, gin.H{"comment": "Delete comment success"})
}
func ArticleList(c *gin.Context) {
	tag := c.Query("tag")
	author := c.Query("author")
	favorited := c.Query("favorited")
	limit := c.Query("limit")
	offset := c.Query("offset")
	articleModels, modelCount, err := models.FindManyArticle(tag, author, limit, offset, favorited)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("articles", errors.New("Invalid param")))
		return
	}
	serializer := serializers.ArticlesSerializer{c, articleModels}
	c.JSON(http.StatusOK, gin.H{"articles": serializer.Response(), "articlesCount": modelCount})
}

func ArticleRetrieve(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "feed" {
		ArticleFeed(c)
		return
	}
	articleModel, err := models.FindOneArticle(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("articles", errors.New("Invalid slug")))
		return
	}
	serializer := serializers.ArticleSerializer{c, articleModel}
	c.JSON(http.StatusOK, gin.H{"article": serializer.Response()})
}


func ArticleCommentList(c *gin.Context) {
	slug := c.Param("slug")
	articleModel, err := models.FindOneArticle(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("comments", errors.New("Invalid slug")))
		return
	}
	err = articleModel.GetComments()
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("comments", errors.New("Database error")))
		return
	}
	serializer := serializers.CommentsSerializer{c, articleModel.Comments}
	c.JSON(http.StatusOK, gin.H{"comments": serializer.Response()})
}
func ArticleFeed(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	myUserModel := c.MustGet("my_user_model").(models.UserModel)
	if myUserModel.ID == 0 {
		c.AbortWithError(http.StatusUnauthorized, errors.New("{error : \"Require auth!\"}"))
		return
	}
	articleUserModel := models.GetArticleUserModel(myUserModel)
	articleModels, modelCount, err := articleUserModel.GetArticleFeed(limit, offset)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("articles", errors.New("Invalid param")))
		return
	}
	serializer := serializers.ArticlesSerializer{c, articleModels}
	c.JSON(http.StatusOK, gin.H{"articles": serializer.Response(), "articlesCount": modelCount})
}

// tag List
func TagList(c *gin.Context) {
	tagModels, err := models.GetAllTags()
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("articles", errors.New("Invalid param")))
		return
	}
	serializer := serializers.TagsSerializer{c, tagModels}
	c.JSON(http.StatusOK, gin.H{"tags": serializer.Response()})
}
