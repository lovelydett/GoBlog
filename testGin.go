package main

import (
    "github.com/gin-gonic/gin"
    "GinBlog/controllers"
)

func main() {
    router:=gin.Default()

    //view path
    router.LoadHTMLGlob("views/*")

    //static path
    router.Static("/static", "./static")

    //routers for GET
    router.GET("/", controllers.HomeGet)
    router.GET("index", controllers.HomeGet)
    router.GET("home", controllers.HomeGet)
    router.GET("article", controllers.ArticleGet)
    router.GET("editArticle", controllers.EditArticleGet)
    router.GET("readArticle", controllers.ReadArticleGet)

    //routers for POST
    router.POST("editArticle", controllers.EditArticlePost)
    router.POST("deleteArticle", controllers.DeleteArticlePost)

    //404 router
    router.NoRoute(controllers.NotFoundGet)

    //start the service
    router.Run(":8080")
}