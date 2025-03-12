package utils

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ListQueryParams struct {
	Filter FilterType `json:"filter"`
	Range  string     `json:"range"`
	Sort   string     `json:"sort"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}

type FilterType struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Sort struct {
	Field string `json:"field"`
	Order string `json:"order"`
}
type Range = []int

type QueryMap struct {
	Filter FilterType `json:"filter" form:"filter"`
	Range  Range      `json:"range" form:"range"`
	Sort   Sort       `json:"sort" form:"sort"`
}

func GetListQueryParams(c *gin.Context) (ListQueryParams, error) {
	var err error

	query := QueryMap{}
	err = c.ShouldBindQuery(&query)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to bind query"})

	}

	return ListQueryParams{
		Filter: query.Filter,
		Range:  c.Query("range"),
		Sort:   c.Query("sort"),
		Limit:  query.Range[1] - query.Range[0],
		Offset: query.Range[0],
	}, nil
}
