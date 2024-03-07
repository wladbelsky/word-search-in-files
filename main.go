package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path/filepath"
	"word-search-in-files/pkg/searcher"
)

func fileSearchEndpoint(c *gin.Context) {
	root, err := filepath.Abs("./examples")
	s := searcher.Searcher{
		FS: os.DirFS(root),
	}
	word := c.Query("word")
	if word == "" {
		c.JSON(400, gin.H{
			"error": "word parameter is required",
		})
		return
	}
	files, err := s.Search(word)
	if err != nil {
		log.Println("search error:", err)
		c.JSON(500, gin.H{
			"error": "search error",
			"files": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"files": files,
	})
}

func main() {

	r := gin.Default()
	r.GET("files/search", fileSearchEndpoint)
	log.Panic(r.Run())
}
