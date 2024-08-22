package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var counter = 0

func main() {
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		counter++
		data := fmt.Sprintf("Server 2 has been accessed %v time(s)", counter)
		ctx.JSON(http.StatusOK, data)
	})
	r.Run(":8002")
}
