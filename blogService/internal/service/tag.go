package service

type CountTagRequest struct {
	State uint8 `from:"state,default=1" binding:"oneof=0 1"`
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
