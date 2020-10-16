package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

//日志系统
var logInf = log.New(os.Stdout, "[INFO]", log.LstdFlags)
var logErr = log.New(os.Stdout, "[Error]", log.LstdFlags | log.Lshortfile)

func HomeGet(c *gin.Context){
	logInf.Println("Entering HomeGet")

	c.HTML(http.StatusOK,"home.html",gin.H{})
	//c.String(http.StatusOK,"ok")

	logInf.Println("Leaving HomeGet")
}