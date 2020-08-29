package model

import (
	"github.com/knight-zlm/blog-service/pkg/app"
)

type Tag struct {
	*Model
	Name  string `json:"name"`
	State string `json:"state"`
}

type TagSwagger struct {
	List  []*Tag
	Pager *app.Pager
}

func (t Tag) TableName() string {
	return "blog_tag"
}
