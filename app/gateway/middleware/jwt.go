package middleware

import (
	"Douyin/util"
	"Douyin/util/resp"
	"errors"
	"github.com/gin-gonic/gin"
)

type Request struct {
	Token string `json:"token" form:"token" binding:"required"`
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 注册和登录接口无需验证
		if c.Request.URL.Path == "/douyin/user/register" || c.Request.URL.Path == "/douyin/user/login" {
			c.Next()
			return
		}

		var request Request
		_ = c.ShouldBind(&request)

		token := request.Token
		// 非视频流接口，token 为空则返回错误
		if c.Request.URL.Path != "/douyin/feed" && token == "" {
			resp.ResponseError(c, resp.CodeNeedLogin, nil)
			c.Abort()
			return
		}

		if token != "" {
			if err := util.ValidateToken(token); err != nil {
				resp.ResponseError(c, resp.CodeError, errors.New("token验证失败"))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
