package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// import "net/http"

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"a": "a",
			"b": []string{"b_todo1", "b_todo2"},
			"e": true,
			"f": false,
			"h": true,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
