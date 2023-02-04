package model

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

type Meta struct {
	TotalPage     int `json:"total_page"`
	CurrentPage   int `json:"current_page"`
	TotalRecord   int `json:"total_record"`
	CurrentRecord int `json:"current_record"`
}
