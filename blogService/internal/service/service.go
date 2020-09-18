package service

import (
	"context"

	"github.com/knight-zlm/blog-service/global"
	"github.com/knight-zlm/blog-service/internal/dao"

	otgorm "github.com/eddycjy/opentracing-gorm"
)

type Service struct {
	ctx context.Context
	dao *dao.Dao
}

func New(c context.Context) Service {
	svc := Service{ctx: c}
	svc.dao = dao.New(otgorm.WithContext(svc.ctx, global.DBEngine))
	return svc
}
