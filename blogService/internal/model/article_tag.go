package model

import "github.com/jinzhu/gorm"

type ArticleTag struct {
	*Model
	TagID     uint32 `json:"tag_id"`
	ArticleID uint32 `json:"article_id"`
}

func (a ArticleTag) TableName() string {
	return "blog_article_tag"
}

func (a ArticleTag) GetByAID(db *gorm.DB) (ArticleTag, error) {
	var articleTag ArticleTag
	err := db.Where("article_id = ? AND is_del = ?", a.ArticleID, 0).First(&articleTag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return articleTag, err
	}
	return articleTag, nil
}

func (a ArticleTag) ListByTID(db *gorm.DB) ([]*ArticleTag, error) {
	var artTags []*ArticleTag
	err := db.Where("tag_id = ? and is_del = ?", a.TagID, 0).Find(&artTags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return artTags, nil
}

func (a ArticleTag) ListByAIDs(db *gorm.DB, artIDs []uint32) ([]*ArticleTag, error) {
	var artTags []*ArticleTag
	err := db.Where("article_id IN (?) AND is_del = ?", artIDs, 0).Find(&artTags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return artTags, nil
}

func (a ArticleTag) Create(db *gorm.DB) (*ArticleTag, error) {
	err := db.Create(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (a ArticleTag) UpdateOne(db *gorm.DB, values interface{}) error {
	err := db.Model(&a).Where("article_id = ? AND is_del = ?", a.ArticleID, 0).Limit(1).Update(values).Error
	if err != nil {
		return err
	}
	return nil
}

func (a ArticleTag) Delete(db *gorm.DB) error {
	err := db.Where("article_id = ? AND is_del = ?", a.ArticleID, a.TagID, 0).Delete(&a).Limit(1).Error
	if err != nil {
		return err
	}
	return nil
}
