package controllers

import (
	"GinBlog/model"
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
	articles := []model.Article{}
	model.GetArticles(0, &articles)
	c.HTML(http.StatusNotFound,"article.html",gin.H{
		"Articles":articles,
	})

	logInf.Println("Leaving NotFoundGet")
}

//editArticle页面的GET请求
func EditArticleGet(c *gin.Context){
	logInf.Println("Entering EditArticleGet for article: ")

	if "yes" == c.Request.FormValue("isNew") {
		logInf.Println("creating a new article")
		//处理编辑新建一个文章的页面
	}else{
		logInf.Println("editing an existing article: ")
		//todo:处理编辑旧文章的页面
	}
	c.HTML(http.StatusOK,"editArticle.html",nil)

	logInf.Println("Leaving EditArticleGet for article: ")
}

type EditArtJson struct{
	IsNew string `json:"isNew"`
	Title string `json:"title"`
	Content string `json:"content"`
}
//editArticle页面的POST请求
func EditArticlePost(c *gin.Context){
	logInf.Println("Entering EditArticlePost: ")

	//需要解析JSON
	artJson := EditArtJson{}
	if nil!=c.BindJSON(&artJson){
		logErr.Println("unable to unmarshal JSON posted from editArticle.html")
		logInf.Println("Leaving EditArticlePost for article: ")
		return
	}

	logInf.Println("JSON: "+artJson.IsNew+", "+artJson.Title+", "+artJson.Content)
	if "yes" == artJson.IsNew {
		logInf.Println("Inserting into db a new article")
		//处理新文章入库的流程
		model.AddArticleByInf(artJson.Title, artJson.Content,1)
		//todo:入库失败校验以及对应的状态码返回
	}else{
		logInf.Println("Updating an existing article in db: ")
		//todo:处理编辑旧文章的页面
	}

	c.HTML(http.StatusOK,"article.html",nil)

	logInf.Println("Leaving EditArticlePost for article: ")
}