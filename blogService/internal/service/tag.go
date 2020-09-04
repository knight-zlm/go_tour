package service

import (
	"github.com/knight-zlm/blog-service/internal/model"
	"github.com/knight-zlm/blog-service/pkg/app"
)

type CountTagRequest struct {
	Name  string `form:"name" binding:"max=10"`
	State uint8  `from:"state,default=1" binding:"oneof=0 1"`
}

type TagListRequest struct {
	Name  string `from:"name" binding:"max=100"`
	State uint8  `from:"state,default=1" binding:"oneof=0 1"`
}

type CreateTagRequest struct {
	Name     string `from:"name" binding:"max=100"`
	CreateBy string `from:"create_by" binding:"required,min=2,max=100"`
	State    uint8  `from:"state,default=1" binding:"oneof=0 1"`
}

type UpdateTagRequest struct {
	ID         uint32 `from:"id" binding:"required,gte=1"`
	Name       string `from:"name" binding:"max=100"`
	State      uint8  `from:"state,default=1" binding:"oneof=0 1"`
	ModifiedBy string `from:"modified_by" binding:"required,min=2,max=100"`
}

type DeleteTagRequest struct {
	ID uint32 `from:"id" binding:"required,gte=1"`
}

func (s *Service) CountTag(param *CountTagRequest) (int, error) {
	return s.dao.CountTag(param.Name, param.State)
}

func (s *Service) GetTagList(param *TagListRequest, pager *app.Pager) ([]*model.Tag, error) {
	return s.dao.GetTagList(param.Name, param.State, pager.Page, pager.PageSize)
}

func (s *Service) UpdateTag(param *UpdateTagRequest) error {
	return s.dao.UpdateTag(param.ID, param.Name, param.State, param.ModifiedBy)
}

func (s *Service) DeleteTag(param *DeleteTagRequest) error {
	return s.dao.DeleteTag(param.ID)
}
