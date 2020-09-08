package service

import "github.com/knight-zlm/blog-service/internal/model"

type Article struct {
	ID            uint32     `json:"id"`
	Title         string     `json:"title"`
	Desc          string     `json:"desc"`
	Content       string     `json:"content"`
	CoverImageUrl string     `json:"cover_image_url"`
	State         uint8      `json:"state"`
	Tag           *model.Tag `json:"tag"`
}

type ArticleRequest struct {
	ID    uint32 `form:"id" binding:"required,gte=1"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

func (s *Service) GetArticle(param *ArticleRequest) (*Article, error) {
	article, err := s.dao.GetArticle(param.ID, param.State)
	if err != nil {
		return nil, err
	}
	articleTag, err := s.dao.GetArticleTagByAID(param.ID)
	if err != nil {
		return nil, err
	}
	tag, err := s.dao.GetTag(articleTag.TagID, model.StateOpen)
	if err != nil {
		return nil, err
	}

	return &Article{
		ID:            article.ID,
		Title:         article.Title,
		Desc:          article.Desc,
		Content:       article.Content,
		CoverImageUrl: article.CoverImageUrl,
		State:         article.State,
		Tag:           &tag,
	}, nil
}
