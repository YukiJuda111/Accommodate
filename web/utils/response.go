package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResponseData(c *gin.Context, errno string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"errno":  errno,
		"errmsg": RecodeText(errno),
		"data":   data,
	})
}
