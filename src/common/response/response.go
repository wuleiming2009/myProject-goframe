package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": data,
	})

}

func failure(c *gin.Context, err error) {
	var bizErr error
	var ok bool
	if bizErr, ok = err.(*BizErr); !ok {
		bizErr = err
	}
	c.JSON(http.StatusOK, bizErr)
}

func Response(c *gin.Context, v interface{}) {
	if err, ok := v.(error); ok {
		failure(c, err)
		return
	}
	success(c, v)
}
