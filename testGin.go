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

    //routers
    router.GET("/", controllers.HomeGet)
    router.GET("index", controllers.HomeGet)
    router.GET("home", controllers.HomeGet)

    //start the service
    router.Run(":8080")
}