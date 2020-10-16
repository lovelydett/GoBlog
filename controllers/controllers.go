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



//home页面的GET请求
func HomeGet(c *gin.Context){
	logInf.Println("Entering HomeGet")

	c.HTML(http.StatusOK,"home.html",nil)
	//c.String(http.StatusOK,"ok")

	logInf.Println("Leaving HomeGet")
}

//404页面的GET请求
func NotFoundGet(c *gin.Context){
	logInf.Println("Entering NotFoundGet for: "+c.Request.URL.String())

	c.HTML(http.StatusNotFound,"404.html",nil)

	logInf.Println("Leaving NotFoundGet")
}

//article页面的GET请求
func ArticleGet(c *gin.Context){
	logInf.Println("Entering ArticleGet")


	//todo:根据页面号进行文章详情的查询
	c.HTML(http.StatusNotFound,"article.html",nil)

	logInf.Println("Leaving NotFoundGet")
}

//editArticle页面的GET请求
func EditArticleGet(c *gin.Context){
	logInf.Println("Entering EditArticleGet for article: ")


	if "yes" == c.Request.FormValue("isNew") {
		logInf.Println("creating a new article")
		//处理新建一个文章的页面
	}else{
		logInf.Println("editing an existing article: ")
		//todo:处理编辑旧文章的页面
	}
	c.HTML(http.StatusOK,"editArticle.html",nil)

	logInf.Println("Leaving EditArticleGet for article: ")
}