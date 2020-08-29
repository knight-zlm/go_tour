package model

import (
	"github.com/knight-zlm/blog-service/pkg/app"
)

type Article struct {
	*Model
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	State         string `json:"state"`
}

type ArticleSwagger struct {
	List  *[]Article
	Pager *app.Pager
}

func (a Article) TableName() string {
	return "blog_article"
}
