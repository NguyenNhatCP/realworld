package serializers

import (
	"crud-rest-api-golang/models"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

type TagSerializer struct {
	C *gin.Context
	models.TagModel
}

type TagsSerializer struct {
	C    *gin.Context
	Tags []models.TagModel
}

func (s *TagSerializer) Response() string {
	return s.TagModel.Tag
}

func (s *TagsSerializer) Response() []string {
	response := []string{}
	for _, tag := range s.Tags {
		serializer := TagSerializer{s.C, tag}
		response = append(response, serializer.Response())
	}
	return response
}

type ArticleUserSerializer struct {
	C *gin.Context
	models.ArticleUserModel
}

func (s *ArticleUserSerializer) Response() ProfileResponse {
	response := ProfileSerializer{s.C, s.ArticleUserModel.UserModel}
	return response.Response()
}

type ArticleSerializer struct {
	C *gin.Context
	models.ArticleModel
}

type ArticleResponse struct {
	ID             uint                  `json:"-"`
	Title          string                `json:"title"`
	Slug           string                `json:"slug"`
	Description    string                `json:"description"`
	Body           string                `json:"body"`
	CreatedAt      string                `json:"createdAt"`
	UpdatedAt      string                `json:"updatedAt"`
	Author         ProfileResponse `json:"author"`
	Tags           []string              `json:"tagList"`
	Favorite       bool                  `json:"favorited"`
	FavoritesCount int64                `json:"favoritesCount"`
}

type ArticlesSerializer struct {
	C        *gin.Context
	Articles []models.ArticleModel
}

func (s *ArticleSerializer) Response() ArticleResponse {
	myUserModel := s.C.MustGet("my_user_model").(models.UserModel)
	authorSerializer := ArticleUserSerializer{s.C, s.Author}
	response := ArticleResponse{
		ID:             s.ID,
		Title:          s.Title,
		Slug:           slug.Make(s.Title),
		Description:    s.Description,
		Body:           s.Body,
		CreatedAt:      s.CreatedAt.UTC().Format("2020-01-02T15:04:05.999Z"),
		UpdatedAt:      s.UpdatedAt.UTC().Format("2020-01-02T15:04:05.999Z"),
		Author:         authorSerializer.Response(),
		Tags:           []string{},
		Favorite:       s.IsFavoriteBy(models.GetArticleUserModel(myUserModel)),
		FavoritesCount: s.FavoritesCount(),
	}
	response.Tags = make([]string, 0)
	for _, tag := range s.Tags {
		serializer := TagSerializer{s.C, tag}
		response.Tags = append(response.Tags, serializer.Response())
	}
	return response
}

func (s *ArticlesSerializer) Response() []ArticleResponse {
	response := []ArticleResponse{}
	for _, article := range s.Articles {
		serializer := ArticleSerializer{s.C, article}
		response = append(response, serializer.Response())
	}
	return response
}

type CommentSerializer struct {
	C *gin.Context
	models.CommentModel
}

type CommentsSerializer struct {
	C        *gin.Context
	Comments []models.CommentModel
}

type CommentResponse struct {
	ID        uint                  `json:"id"`
	Body      string                `json:"body"`
	CreatedAt string                `json:"createdAt"`
	UpdatedAt string                `json:"updatedAt"`
	Author    ProfileResponse `json:"author"`
}

func (s *CommentSerializer) Response() CommentResponse {
	authorSerializer := ArticleUserSerializer{s.C, s.Author}
	response := CommentResponse{
		ID:        s.ID,
		Body:      s.Body,
		CreatedAt: s.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt: s.UpdatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		Author:    authorSerializer.Response(),
	}
	return response
}

func (s *CommentsSerializer) Response() []CommentResponse {
	response := []CommentResponse{}
	for _, comment := range s.Comments {
		serializer := CommentSerializer{s.C, comment}
		response = append(response, serializer.Response())
	}
	return response
}