package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	port := gin.Default()
	port.Run("localhost:8000")
}
