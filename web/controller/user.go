package controller

import (
	"RentHouse/web/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSession(c *gin.Context) {
	errno := utils.RECODE_SESSIONERR
	errmsg := utils.RecodeText(errno)
	c.JSON(http.StatusOK, gin.H{
		"errno":  errno,
		"errmsg": errmsg,
	})
}
