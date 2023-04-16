package model

type Pagination struct {
	Limit  int    `json:"limit" query:"limit"`
	Page   int    `json:"page" query:"page"`
	Sort   string `json:"sort" query:"sort"`
	Search string `json:"search" query:"search"`
}

type Meta struct {
	TotalPage     int `json:"total_page"`
	CurrentPage   int `json:"current_page"`
	TotalRecord   int `json:"total_record"`
	CurrentRecord int `json:"current_record"`
}
