package api

type PaginationQueryParams struct {
	Limit int32 `form:"limit" binding:"required,min=5,max=10"`
	Skip  int32 `form:"skip"`
}

type PathIDParam struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}
