package utils

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
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
	filter := FilterType{}
	filterQuery := c.Query("filter")
	if filterQuery != "" {
		if err := json.Unmarshal([]byte(filterQuery), &filter); err != nil {
			return ListQueryParams{}, err
		}
	}
	if filter.Start == "" {
		filter.Start = c.Query("filter.start")
	}
	if filter.End == "" {
		filter.End = c.Query("filter.end")
	}

	queryRange, err := parseRangeQuery(c)
	if err != nil {
		return ListQueryParams{}, err
	}

	return ListQueryParams{
		Filter: filter,
		Range:  c.Query("range"),
		Sort:   c.Query("sort"),
		Limit:  queryRange[1] - queryRange[0],
		Offset: queryRange[0],
	}, nil
}

func parseRangeQuery(c *gin.Context) (Range, error) {
	rangeValues := c.QueryArray("range")
	if len(rangeValues) == 0 {
		return nil, errors.New("range query is required")
	}

	if len(rangeValues) == 1 {
		var queryRange Range
		if err := json.Unmarshal([]byte(rangeValues[0]), &queryRange); err != nil {
			return nil, err
		}
		if len(queryRange) < 2 {
			return nil, errors.New("range query must include start and end")
		}
		return queryRange, nil
	}

	start, err := strconv.Atoi(rangeValues[0])
	if err != nil {
		return nil, err
	}
	end, err := strconv.Atoi(rangeValues[1])
	if err != nil {
		return nil, err
	}
	return Range{start, end}, nil
}
