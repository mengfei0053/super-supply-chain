package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func SetContentRange(c *gin.Context, input int64) {
	c.Header("Content-Range", strconv.FormatInt(input, 10))
}
