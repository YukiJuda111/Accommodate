package middleware

import (
	"RentHouse/web/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LoginFilter 登录过滤器中间件
func LoginFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取session
		session := sessions.Default(c)
		userName := session.Get("userName")
		if userName == nil {
			c.JSON(http.StatusOK, gin.H{
				"errno":  utils.RECODE_SESSIONERR,
				"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
