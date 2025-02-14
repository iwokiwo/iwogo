package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func SetTimeZone(timezone string) gin.HandlerFunc {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		panic(fmt.Sprintf("Invalid timezone: %s", err))
	}

	return func(c *gin.Context) {
		// Set the timezone for the current request
		c.Set("timezone", loc)
		c.Next()
	}
}
