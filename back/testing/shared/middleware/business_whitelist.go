package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var allowedBusinessIDs = map[uint]bool{
	4:  true,
	7:  true,
	26: true,
	33: true,
}

func BusinessWhitelist() gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Query("business_id")
		if param == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "business_id query parameter is required"})
			c.Abort()
			return
		}

		id, err := strconv.ParseUint(param, 10, 64)
		if err != nil || id == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid business_id"})
			c.Abort()
			return
		}

		if !allowedBusinessIDs[uint(id)] {
			c.JSON(http.StatusForbidden, gin.H{
				"error":       "Business not authorized for testing",
				"allowed_ids": []uint{4, 7, 26, 33},
			})
			c.Abort()
			return
		}

		c.Set("testing_business_id", uint(id))
		c.Next()
	}
}

func GetAllowedBusinessIDs() []uint {
	return []uint{4, 7, 26, 33}
}
