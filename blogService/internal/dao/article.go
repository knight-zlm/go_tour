package dao

import "github.com/knight-zlm/blog-service/internal/model"

type Article struct {
	ID            uint32 `json:"id"`
	TagID         uint32 `json:"tag_id"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	CreateBy      string `json:"create_by"`
	ModifiedBy    string `json:"modified_by"`
	State         uint8  `json:"state"`
}

func (d *Dao) CreateArticle(param *Article) (*model.Article, error) {
	art := model.Article{
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		State:         param.State,
		Model:         &model.Model{CreatedBy: param.CreateBy},
	}
	return &art, nil
}

func (d *Dao) UpdateArticle(param *Article) error {
	art := model.Article{
		Model: &model.Model{ID: param.ID},
	}
	values := map[string]interface{}{
		"modified_by": param.ModifiedBy,
		"state":       param.State,
	}
	if param.Title != "" {
		values["title"] = param.Title
	}
	if param.CoverImageUrl != "" {
		values["cover_image_url"] = param.CoverImageUrl
	}
	if param.Desc != "" {
		values["desc"] = param.Desc
	}
	if param.Content != "" {
		values["content"] = param.Content
	}
	return art.Update(d.engine, values)
}

func (d *Dao) GetArticle(id uint32, state uint8) (*model.Article, error) {
	art := model.Article{
		State: state,
		Model: &model.Model{ID: id},
	}
	return art.Get(d.engine)
}

func (d *Dao) DeleteArticle(id uint32) error {
	art := model.Article{
		Model: &model.Model{ID: id},
	}
	return art.Delete(d.engine)
}

func (d *Dao) CountArticleListByTagID(id uint32, state uint8) (int, error) {
	art := model.Article{
		State: state,
	}
	return art.CountByTagID(d.engine, id)
}

func (d *Dao) GetArticleListByTagID(id uint32, state uint8, page, pageSize int) ([]*model.ArticleRow, error) {
	art := model.Article{
		State: state,
	}
	return art.ListByTagID(d.engine, id, page, pageSize)
}
