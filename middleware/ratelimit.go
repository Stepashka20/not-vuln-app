package middleware

import (
	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"time"
)

var mw gin.HandlerFunc

func GetMW() gin.HandlerFunc {
	rateStore := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Second,
		Limit: 2,
	})
	mw = ratelimit.RateLimiter(rateStore, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	return mw
}

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.HTML(429, "error.html", gin.H{"message": "Слишком много запросов, попробуйте позже"})
}
