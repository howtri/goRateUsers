package main

import (
	"github.com/gin-gonic/gin"
	"github.com/howtri/goRateUsers/handlers"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		c.File("./ui/index.html")
	})

	r.POST("/user/add", handlers.AddUserHandler)
	r.POST("/user/verify", handlers.VerifyUserHandler)

	err := r.Run(":3001")
	if err != nil {
		panic(err)
	}

}
