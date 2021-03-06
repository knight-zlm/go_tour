package routers

import (
	"net/http"
	"time"

	_ "github.com/knight-zlm/blog-service/docs"
	"github.com/knight-zlm/blog-service/global"
	"github.com/knight-zlm/blog-service/internal/middleware"
	"github.com/knight-zlm/blog-service/internal/routers/api"
	v1 "github.com/knight-zlm/blog-service/internal/routers/api/v1"
	"github.com/knight-zlm/blog-service/internal/routers/upload"
	"github.com/knight-zlm/blog-service/pkg/limiter"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var methodLimiter = limiter.NewMethodLimiter().AddBucket(limiter.LimiterBucketRule{
	Key:          "/auth",
	FillInterval: time.Second,
	Capacity:     10,
	Quantum:      10,
})

func NewRouter() *gin.Engine {
	r := gin.New()
	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger(), gin.Recovery())
	} else {
		r.Use(middleware.AccessLog(), middleware.Recovery())
	}
	r.Use(middleware.Tracing())
	r.Use(middleware.RateLimiter(methodLimiter))
	r.Use(middleware.ContextTimeOut(global.AppSetting.DefaultContextTimeout))
	r.Use(middleware.Translations())
	// url:= ginSwagger.URL("http://127.0.0.1:8008/swagger/doc.json")
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload/file", upload.UploadFile)
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadServerUrl))
	r.GET("/auth", api.GetAuth)

	article := v1.NewArticle()
	tag := v1.NewTag()
	apiv1 := r.Group("/api/v1")
	// 鉴权中间件
	//apiv1.Use(middleware.JWT())
	{
		apiv1.POST("/tags", tag.Create)
		apiv1.DELETE("/tags/:id", tag.Delete)
		apiv1.PUT("/tags/:id", tag.Update)
		apiv1.PATCH("/tags/:id/state", tag.Update)
		apiv1.GET("/tags", tag.List)

		apiv1.POST("/articles", article.Create)
		apiv1.DELETE("/articles/:id", article.Delete)
		apiv1.PUT("/articles/:id", article.Update)
		apiv1.PATCH("/articles/:id/state", article.Update)
		apiv1.GET("/articles/:id", article.Get)
		apiv1.GET("/articles", article.List)
	}
	return r
}
