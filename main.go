package main

import (
	"github.com/gin-gonic/gin"
	"github.com/howtri/goRate/handlers"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		c.File("./ui/index.html")
	})

	r.POST("/skill/add", handlers.AddSkillHandler)
	r.POST("/skill/search", handlers.SearchSkillsHandler)
	r.POST("/skill/rank", handlers.RankSkillHandler)
	r.GET("/skill/:id", handlers.GetSkillHandler)

	err := r.Run(":3000")
	if err != nil {
		panic(err)
	}

}
