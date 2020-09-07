package dao

import "github.com/knight-zlm/blog-service/internal/model"

func (d *Dao) GetArticleTagByAID(artID uint32) (model.ArticleTag, error) {
	artTag := model.ArticleTag{ArticleID: artID}
	return artTag.GetByAID(d.engine)
}

func (d *Dao) GetArticleTagListByAIDs(artIDs []uint32) ([]*model.ArticleTag, error) {
	artTag := model.ArticleTag{}
	return artTag.ListByAIDs(d.engine, artIDs)
}

func (d *Dao) GetArticleTagListByTID(tagID uint32) ([]*model.ArticleTag, error) {
	artTag := model.ArticleTag{TagID: tagID}
	return artTag.ListByTID(d.engine)
}

func (d *Dao) CreateArticleTag(artID, tagID uint32, createBy string) (*model.ArticleTag, error) {
	artTag := model.ArticleTag{ArticleID: artID, TagID: tagID, Model: &model.Model{CreatedBy: createBy}}
	return artTag.Create(d.engine)
}

func (d *Dao) UpdateArticleTag(artID, tagID uint32, modifiedBy string) error {
	artTag := model.ArticleTag{ArticleID: artID, TagID: tagID, Model: &model.Model{ModifiedBy: modifiedBy}}
	values := map[string]interface{}{
		"article_id":  artID,
		"tag_id":      tagID,
		"modified_by": modifiedBy,
	}
	return artTag.UpdateOne(d.engine, values)
}

func (d *Dao) DeleteArticleTag(artID uint32) error {
	artTag := model.ArticleTag{ArticleID: artID}
	return artTag.DeleteOne(d.engine)
}
