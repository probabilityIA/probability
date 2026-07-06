package ratelimit

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type KeyFunc func(c *gin.Context) string

func Gin(l Limiter, keyFunc KeyFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if l == nil || keyFunc == nil {
			c.Next()
			return
		}
		key := keyFunc(c)
		if key == "" {
			c.Next()
			return
		}

		d := l.Check(c.Request.Context(), key)
		if d.Allowed {
			c.Next()
			return
		}

		retry := int(d.RetryAfter.Seconds())
		if retry < 1 {
			retry = 1
		}
		c.Header("Retry-After", strconv.Itoa(retry))
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"error":       "rate_limited",
			"blacklisted": d.Blacklisted,
			"retry_after": retry,
		})
	}
}

func ByHeader(header, prefix string) KeyFunc {
	return func(c *gin.Context) string {
		if v := c.GetHeader(header); v != "" {
			return prefix + ":" + v
		}
		return ""
	}
}

func ByParam(param, prefix string) KeyFunc {
	return func(c *gin.Context) string {
		if v := c.Param(param); v != "" {
			return prefix + ":" + v
		}
		return ""
	}
}

func ByClientIP(prefix string) KeyFunc {
	return func(c *gin.Context) string {
		return prefix + ":" + c.ClientIP()
	}
}

func FirstNonEmpty(funcs ...KeyFunc) KeyFunc {
	return func(c *gin.Context) string {
		for _, f := range funcs {
			if k := f(c); k != "" {
				return k
			}
		}
		return ""
	}
}
