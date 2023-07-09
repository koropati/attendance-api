package model

type Pagination struct {
	Limit  int    `json:"limit" query:"limit" form:"limit"`
	Page   int    `json:"page" query:"page" form:"page"`
	Sort   string `json:"sort" query:"sort" form:"sort"`
	Search string `json:"search" query:"search" form:"search"`
}

type Meta struct {
	TotalPage     int `json:"total_page" query:"total_page" form:"total_page"`
	CurrentPage   int `json:"current_page" query:"current_page" form:"current_page"`
	TotalRecord   int `json:"total_record" query:"total_record" form:"total_record"`
	CurrentRecord int `json:"current_record" query:"current_record" form:"current_record"`
}
