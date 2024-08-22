package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var counter = 0

func main() {
	r := gin.Default()
	counter++
	r.GET("/", func(ctx *gin.Context) {
		data := fmt.Sprintf("Server 1 has been accessed %v time(s) ", counter)
		ctx.JSON(http.StatusOK, data)
	})

	r.Run(":8000")
}
