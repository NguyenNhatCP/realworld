package models

import (
	"crud-rest-api-golang/database"
	"strconv"

	"github.com/golang/glog"
	"gorm.io/gorm"
)

type ArticleModel struct {
	gorm.Model
	Slug        string `gorm:"unique_index"`
	Title       string
	Description string `gorm:"size:2048"`
	Body        string `gorm:"size:2048"`
	Author      ArticleUserModel
	AuthorID    uint
	Tags        []TagModel     `gorm:"many2many:article_tags;"`
	Comments    []CommentModel `gorm:"ForeignKey:ArticleID"`
}

type ArticleUserModel struct {
	gorm.Model
	UserModel      UserModel
	UserModelID    uint
	ArticleModels  []ArticleModel  `gorm:"ForeignKey:AuthorID"`
	FavoriteModels []FavoriteModel `gorm:"ForeignKey:FavoriteByID"`
}

type FavoriteModel struct {
	gorm.Model
	Favorite     ArticleModel
	FavoriteID   uint
	FavoriteBy   ArticleUserModel
	FavoriteByID uint
}

type TagModel struct {
	gorm.Model
	Tag           string         `gorm:"unique_index"`
	ArticleModels []ArticleModel `gorm:"many2many:article_tags;"`
}

type CommentModel struct {
	gorm.Model
	Article   ArticleModel
	ArticleID uint
	Author    ArticleUserModel
	AuthorID  uint
	Body      string `gorm:"size:2048"`
}

func GetArticleUserModel(userModel UserModel) ArticleUserModel {
	var articleUserModel ArticleUserModel
	if userModel.ID == 0 {
		return articleUserModel
	}
	db := database.GetDB()
	db.Where(&ArticleUserModel{
		UserModelID: userModel.ID,
	}).FirstOrCreate(&articleUserModel)
	articleUserModel.UserModel = userModel
	return articleUserModel
}

func (model *ArticleModel) SetTags(tags []string) error {
	db := database.GetDB()
	var tagList []TagModel
	for _, tag := range tags {
		var tagModel TagModel
		err := db.FirstOrCreate(&tagModel, TagModel{Tag: tag}).Error
		if err != nil {
			return err
		}
		tagList = append(tagList, tagModel)
	}
	model.Tags = tagList
	return nil
}

func (article ArticleModel) FavoritesCount() int64 {
	db := database.GetDB()
	var count int64
	db.Model(&FavoriteModel{}).Where(FavoriteModel{
		FavoriteID: article.ID,
	}).Count(&count)
	return count
}

func (article ArticleModel) IsFavoriteBy(user ArticleUserModel) bool {
	db := database.GetDB()
	var favorite FavoriteModel
	db.Where(FavoriteModel{
		FavoriteID:   article.ID,
		FavoriteByID: user.ID,
	}).First(&favorite)
	return favorite.ID != 0
}

func (article ArticleModel) Update(data interface{}) error {
	db := database.GetDB()
	err := db.Save(data).Error
	return err
}

func FindOneArticle(s string) (ArticleModel, error) {
	db := database.GetDB()
	var model ArticleModel
	tx := db.Begin()
	tx.Model(&model).Where(ArticleModel{Slug: s}).Preload("Author.UserModel").Preload("Tags").Find(&model)
	// select * from article_model
	// select * from Author.UserModel where UserModelID IN (1,2,3) // belong to
	// select * from Tags where article_id IN (1,2,3) // has many

	err := tx.Commit().Error
	glog.Info("Starting transaction..." + tx.Statement.Table)
	return model, err
}

func DeleteArticleModel(condition interface{}) error {
	db := database.GetDB()
	err := db.Where(condition).Delete(&ArticleModel{}).Error
	return err
}

func (article ArticleModel) FavoriteBy(user ArticleUserModel) error {
	db := database.GetDB()
	var favorite FavoriteModel
	err := db.FirstOrCreate(&favorite, &FavoriteModel{
		FavoriteID:   article.ID,
		FavoriteByID: user.ID,
	}).Error
	return err
}

func (article ArticleModel) UnFavoriteBy(user ArticleUserModel) error {
	db := database.GetDB()
	err := db.Where(FavoriteModel{
		FavoriteID:   article.ID,
		FavoriteByID: user.ID,
	}).Delete(&FavoriteModel{}).Error
	return err
}
func DeleteCommentModel(condition interface{}) error {
	db := database.GetDB()
	err := db.Where(condition).Delete(&CommentModel{}).Error
	return err
}


func (model *ArticleModel) GetComments() error {
	db := database.GetDB()
	tx := db.Begin()
	tx.Preload("Comments").Find(&model)
	for i, _ := range model.Comments {
		tx.Preload("Author.UserModel").Find(&model.Comments[i])
	}
	err := tx.Commit().Error
	return  err
}
func FindManyArticle(tag, author, limit, offset, favorited string) ([]ArticleModel, int64, error) {
	db := database.GetDB()
	var model []ArticleModel
	var count int64

	offset_int, err := strconv.Atoi(offset)
	if err != nil {
		offset_int = 0
	}

	limit_int, err := strconv.Atoi(limit)
	if err != nil {
		limit_int = 20
	}

	tx := db.Begin()
	if tag != "" {
		var tagModel TagModel
		tx.Where(TagModel{Tag: tag}).First(&tagModel)
		if tagModel.ID != 0 {
			tx.Model(&tagModel).
			Preload("Author.UserModel").
			Preload("Tags").
			Offset(offset_int).
			Limit(limit_int).
			Order("created_at desc").
			Association("ArticleModels").
			Find(&model)
			//SELECT `article_models`.`id`,`article_models`.`created_at`,`article_models`.`updated_at`,`article_models`.`deleted_at`,`article_models`.`tag` FROM `article_models` JOIN `article_tags` ON `article_tags`.`article_model_id` = `article_models`.`id` AND `article_tags`.`tag_model_id` = 3 WHERE `article_models`.`deleted_at` IS NULL ORDER BY created_at desc LIMIT 20
			// tx.Model(&tagModel).Offset(offset_int).Limit(limit_int).Related(&model, "ArticleModels")
			count = tx.Model(&tagModel).Association("ArticleModels").Count()
		}
	} else if author != "" {
		var userModel UserModel
		tx.Where(UserModel{Username: author}).First(&userModel)
		articleUserModel := GetArticleUserModel(userModel)

		if articleUserModel.ID != 0 {
			tx.Debug().Model(&articleUserModel).
			Preload("Author.UserModel").
			Preload("Tags").
			Offset(offset_int).
			Limit(limit_int).
			Order("created_at desc").
			Association("ArticleModels").
			Find(&model)
			// tx.Preload("ArticleModels").Find(&model).Offset(offset_int).Limit(limit_int)
			// tx.Model(&articleUserModel).Offset(offset_int).Limit(limit_int).Related(&model, "ArticleModels")
			count = tx.Model(&articleUserModel).Association("ArticleModels").Count()
		}
	} else if favorited != "" {
		var userModel UserModel
		tx.Where(UserModel{Username: favorited}).First(&userModel)
		articleUserModel := GetArticleUserModel(userModel)
		if articleUserModel.ID != 0 {
			var favoriteModels []FavoriteModel
			tx.Where(&FavoriteModel{
				FavoriteByID: articleUserModel.ID,
			}).Offset(offset_int).Limit(limit_int).Find(&favoriteModels)

			for _, favorite := range favoriteModels {
				var articleModel ArticleModel
				// tx.Model(&favorite).Related(&model, "Favorite")
				tx.Model(&favorite).
				// Preload("Favorite").
				Preload("Author.UserModel").
				Preload("Tags").
				Offset(offset_int).
				Limit(limit_int).
				Order("created_at desc").
				Association("Favorite").
				Find(&model)
				if(articleModel.ID != 0) {
					model = append(model, articleModel)
				}
			}
			count = tx.Model(&articleUserModel).Association("FavoriteModels").Count()
			// SELECT * FROM `favorite_models` WHERE `favorite_models`.`favorite_id` = 4 AND `favorite_models`.`deleted_at` IS NULL ORDER BY `favorite_models`.`id` LIMIT 1
		}
	} else {
		db.Model(&model).Count(&count)
		db.Offset(offset_int).Limit(limit_int).Find(&model)
		for i, _ := range model {
		tx.Model(&model[i]).Preload("Author.UserModel").Preload("Tags").Find(&model)
		}
	}
	err = tx.Commit().Error
	return model, count, err
}
func (model * ArticleUserModel) GetArticleFeed(limit, offset string) ([]ArticleModel, int64, error) {
	db := database.GetDB()
	var models []ArticleModel
	var count int64

	offset_int, err := strconv.Atoi(offset)
	if err != nil {
		offset_int = 0
	}
	limit_int, err := strconv.Atoi(limit)
	if err != nil {
		limit_int = 20
	}

	tx := db.Begin()
	followings := model.UserModel.GetFollowings()
	var articleUserModels []uint
	for _, following := range followings {
		articleUserModel := GetArticleUserModel(following)
		articleUserModels = append(articleUserModels, articleUserModel.ID)
	}

	tx.Debug().Where("author_id in (?)", articleUserModels).Order("updated_at desc").Offset(offset_int).Limit(limit_int).Find(&models)

	for i, _ := range models {
		tx.Debug().Model(&models[i]).Preload("Author.UserModel").Preload("Tags").Find(&models)
		count++;
	}
	err = tx.Commit().Error
	return models, count, err
}
// get all tag
func GetAllTags() ([]TagModel, error) {
	db := database.GetDB()
	var models []TagModel
	err := db.Find(&models).Error
	return models, err
}
