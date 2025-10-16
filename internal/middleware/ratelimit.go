package middleware

import (
	"net/http"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

// errorHandler - то, что приходит, когда лимит превышен
func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
		"error": "Too many requests",
		"code":  "rate_limit_exceeded",
	})

}

func NewRateLimiter(rdb *redis.Client) gin.HandlerFunc {
	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: rdb,
		Rate:        time.Minute,
		Limit:       4,
	})

	return ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})
}
