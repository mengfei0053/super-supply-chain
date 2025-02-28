package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"log"
	"net/http"
)

type ListQueryParams struct {
	Filter string `json:"filter"`
	Range  string `json:"range"`
	Sort   string `json:"sort"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

func GetListQueryParams(c *gin.Context) (ListQueryParams, error) {
	var sliceRange []int
	range_ := c.Query("range")
	
	err := json.Unmarshal([]byte(range_), &sliceRange)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal range"})
		return ListQueryParams{}, err
	}

	limit := sliceRange[1] - sliceRange[0] + 1
	offset := sliceRange[0]

	return ListQueryParams{
		Filter: c.Query("filter"),
		Range:  c.Query("range"),
		Sort:   c.Query("sort"),
		Limit:  limit,
		Offset: offset,
	}, nil
}
