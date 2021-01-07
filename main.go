package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"redisManger/src/handlers"
)

func main() {
	r := gin.Default()
	r.GET("/redis", handlers.GetTop20)

	log.Fatal(r.Run(":8899"))

}
