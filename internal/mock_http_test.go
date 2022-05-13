package internal

import (
	"log"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestParam(t *testing.T) {
	engion := gin.New()
	engion.Any("/restaurant/:id/order/:orderNo", func(c *gin.Context) {
		log.Println(c.Params)

	})
}
