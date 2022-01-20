package validator

import (
	"crud-rest-api-golang/common"
	"crud-rest-api-golang/models"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)
type ArticleModelValidator struct {
	Article struct {
		Title       string   `form:"title" json:"title" binding:"required,min=4"`
		Description string   `form:"description" json:"description" binding:"max=2048"`
		Body        string   `form:"body" json:"body" binding:"max=2048"`
		Tags        []string `form:"tagList" json:"tagList"`
	} `json:"article"`
	ArticleModel models.ArticleModel `json:"-"`
}

func NewArticleModelValidator() ArticleModelValidator {
	return ArticleModelValidator{}
}

func NewArticleModelValidatorFillWith(articleModel models.ArticleModel) ArticleModelValidator {
	articleModelValidator := NewArticleModelValidator()
	articleModelValidator.Article.Title = articleModel.Title
	articleModelValidator.ArticleModel.CreatedAt = articleModel.CreatedAt
	articleModelValidator.Article.Description = articleModel.Description
	articleModelValidator.Article.Body = articleModel.Body
	for _, tagModel := range articleModel.Tags {
		articleModelValidator.Article.Tags = append(articleModelValidator.Article.Tags, tagModel.Tag)
	}
	return articleModelValidator
}

func (s *ArticleModelValidator) Bind(c *gin.Context) error {
	myUserModel := c.MustGet("my_user_model").(models.UserModel)

	err := common.Bind(c, s)
	if err != nil {
		return err
	}
	s.ArticleModel.Slug = slug.Make(s.Article.Title)
	s.ArticleModel.Title = s.Article.Title
	s.ArticleModel.Description = s.Article.Description
	s.ArticleModel.Body = s.Article.Body
	s.ArticleModel.Author = models.GetArticleUserModel(myUserModel)
	s.ArticleModel.SetTags(s.Article.Tags)
	return nil
}

type CommentModelValidator struct {
	Comment struct {
		Body string `json:"body" binding:"max=2048"`
	} `json:"comment"`
	CommentModel models.CommentModel `json:"-"`
}

func NewCommentModelValidator() CommentModelValidator {
	return CommentModelValidator{}
}

func (s *CommentModelValidator) Bind(c *gin.Context) error {
	myUserModel := c.MustGet("my_user_model").(models.UserModel)

	err := common.Bind(c, s)
	if err != nil {
		return err
	}
	s.CommentModel.Body = s.Comment.Body
	s.CommentModel.Author = models.GetArticleUserModel(myUserModel)
	return nil
}
