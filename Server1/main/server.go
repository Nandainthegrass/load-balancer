package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	counter int
	mutex   sync.Mutex
)

func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		mutex.Lock()         // Lock the mutex before accessing the shared resource
		defer mutex.Unlock() // Ensure the mutex is unlocked when the function completes

		counter++
		data := fmt.Sprintf("Server 1 has been accessed %v time(s)", counter)
		time.Sleep(3 * time.Second)
		ctx.JSON(http.StatusOK, data)
	})

	r.Run(":8000")
}
