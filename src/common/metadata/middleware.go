package metadata

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"myProject/common/response"
	"myProject/common/token"
	"myProject/lib/log"
)

func LogWithReqId(c *gin.Context) {
	c.Request = c.Request.WithContext(log.NewContext(c.Request.Context(), GenReqId()))
}

// 将User-Agent放到ctx中
func ParseUserAgent(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, KeyUserAgent, c.Request.UserAgent())
	c.Request = c.Request.WithContext(ctx)
}

func UserAgentFromCtx(ctx context.Context) string {
	if userAgent, ok := ctx.Value(KeyUserAgent).(string); ok {
		return userAgent
	}
	return ""
}

func WithReqId(c *gin.Context) {
	ctx := c.Request.Context()
	reqId := GenReqId()
	ctx = NewContextWithReqId(ctx, reqId)
	c.Request = c.Request.WithContext(ctx)
	c.Writer.Header().Set(headerKeyReqId, reqId)
}

func Auth(c *gin.Context) {
	ctx := c.Request.Context()

	userToken := c.Request.Header.Get("token")
	userId, err := token.GetUserIdByToken(ctx, userToken)
	if err != nil {
		response.Response(c, response.ErrRelogin)
		c.Abort()
		return
	}
	ctx = NewContextWithUserId(ctx, userId)
	c.Request = c.Request.WithContext(ctx)
}

// 跨域问题
func Cors() gin.HandlerFunc {
	// CORS for https://foo.com and https://github.com origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodHead, http.MethodOptions},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "X-Reqid"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "https://github.com"
		//},
		MaxAge: 12 * time.Hour,
	})
}
