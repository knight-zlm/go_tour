package service

import (
	"github.com/knight-zlm/blog-service/internal/dao"
	"github.com/knight-zlm/blog-service/internal/model"
	"github.com/knight-zlm/blog-service/pkg/app"
)

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

type ArticleListRequest struct {
	ID    uint32 `form:"id" binding:"required,gte=1"`
	TagID uint32 `form:"tag_id" binding:"required,gte=1"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type CreateArticleRequest struct {
	ID            uint32 `json:"id"`
	TagID         uint32 `json:"tag_id"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	State         uint8  `json:"state"`
	CreatedBy     string `json:"created_by"`
}

type UpdateArticleRequest struct {
	ID            uint32 `json:"id"`
	TagID         uint32 `json:"tag_id"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	State         uint8  `json:"state"`
	ModifiedBy    string `json:"modified_by"`
}

type DeleteArticleRequest struct {
	ID         uint32 `json:"id"`
	ModifiedBy string `json:"modified_by"`
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

func (s *Service) GetArticleList(param *ArticleListRequest, pager *app.Pager) ([]*Article, error) {
	artCount, err := s.dao.CountArticleListByTagID(param.TagID, param.State)
	if err != nil {
		return nil, err
	}
	articles, err := s.dao.GetArticleListByTagID(param.TagID, param.State, pager.Page, pager.PageSize)
	if err != nil {
		return nil, err
	}
	artList := make([]*Article, 0, artCount)
	for _, article := range articles {
		artList = append(artList, &Article{
			ID:            article.ArticleID,
			Title:         article.ArticleTitle,
			Desc:          article.ArticleDesc,
			Content:       article.Content,
			CoverImageUrl: article.CoverImageUrl,
			Tag: &model.Tag{
				Model: &model.Model{ID: article.TagID},
				Name:  article.TagName,
			},
		})
	}
	return artList, nil
}

func (s *Service) CreateArticle(param *CreateArticleRequest) error {
	article, err := s.dao.CreateArticle(&dao.Article{
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		CreateBy:      param.CreatedBy,
		State:         param.State,
	})
	if err != nil {
		return err
	}

	_, err = s.dao.CreateArticleTag(article.ID, param.TagID, param.CreatedBy)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateArticle(param *UpdateArticleRequest) error {
	err := s.dao.UpdateArticle(&dao.Article{
		ID:            param.ID,
		TagID:         param.TagID,
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		ModifiedBy:    param.ModifiedBy,
		State:         param.State,
	})
	if err != nil {
		return err
	}

	err = s.dao.UpdateArticleTag(param.ID, param.TagID, param.ModifiedBy)
	if err != nil {
		return err
	}

	return err
}

func (s *Service) DeleteArticle(param *DeleteArticleRequest) error {
	err := s.dao.DeleteArticle(param.ID)
	if err != nil {
		return err
	}

	err = s.dao.DeleteArticleTag(param.ID)
	if err != nil {
		return err
	}

	return nil
}
