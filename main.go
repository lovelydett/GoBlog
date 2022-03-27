package main

import (
	"GinBlog/controllers"
	. "GinBlog/utils"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

/* logging */
var logInf = log.New(os.Stdout, "[INFO]", log.LstdFlags)
var logErr = log.New(os.Stdout, "[Error]", log.LstdFlags|log.Lshortfile)

// do init works here
func Init() error {
	// 1. check and create the video path
	working_path, _ := os.Getwd()
	video_path := working_path + "/static/video"
	if !IsPathExist(video_path) {
		err := os.MkdirAll(video_path, 777)
		err = os.Chmod(video_path, 777)
		if err != nil {
			logErr.Println("Unable to create video path:", video_path)
			return err
		}
	}
	return nil
}

func main() {

	err := Init()
	if err != nil {
		logErr.Println("Unexpected error during init, program exits")
		os.Exit(-1)
	}

	router := gin.Default()

	// view path
	router.LoadHTMLGlob("new_views/*")

	// static path
	router.Static("/static", "./static")

	// routers for GET
	router.GET("/", controllers.HomeGet)
	router.GET("/article", controllers.ArticleGet)
	router.GET("/readArticle", controllers.ReadArticleGet)
	router.GET("/editArticle", controllers.EditArticleGet)

	// routers for POST
	router.POST("/editArticle", controllers.EditArticlePost)

	// start the service
	router.Run(":8080")
}
