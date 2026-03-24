package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"test.com/project-api/api/rpc"
	common "test.com/project-common"
	"test.com/project-common/errs"
	auth "test.com/project-grpc/auth"
)

var ignores = []string{
	"project/login/register",
	"project/login",
	"project/login/getCaptcha",
	"project/organization",
	"project/auth/apply"}

func Auth() func(*gin.Context) {
	return func(c *gin.Context) {
		result := &common.Result{}
		uri := c.Request.RequestURI
		for _, v := range ignores {
			if strings.Contains(uri, v) {
				c.Next()
				return
			}
		}
		memberId := c.GetInt64("memberId")
		msg := &auth.AuthReqMessage{
			MemberId: memberId,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		response, err := rpc.AuthServiceClient.AuthNodesByMemberId(ctx, msg)
		if err != nil {
			code, msg := errs.ParseGrpcError(err)
			c.JSON(http.StatusOK, result.Fail(code, msg))
			c.Abort()
			return
		}

		for _, v := range response.List {
			if strings.Contains(uri, v) {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusOK, result.Fail(403, "无权限操作"))
		c.Abort()
		return
	}
}
