package pagination

import (
	"attendance-api/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GeneratePaginationFromRequest(c *gin.Context) model.Pagination {
	limit := 2
	page := 1
	sort := "created_at asc"
	query := c.Request.URL.Query()
	for key, value := range query {
		queryValue := value[len(value)-1]
		if key == "limit" {
			limit, _ = strconv.Atoi(queryValue)
		} else if key == "page" {
			page, _ = strconv.Atoi(queryValue)
		} else {
			sort = queryValue
		}
	}
	return model.Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}
}
