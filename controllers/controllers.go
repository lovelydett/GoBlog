package controllers

import (
	"GinBlog/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
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
	c.HTML(http.StatusOK,"article.html",gin.H{
		"Articles":articles,
	})

	logInf.Println("Leaving ArticleGet")
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
		c.JSON(http.StatusOK,gin.H{"status":1, "message":"Unable to post an article", "data":nil})
		return
	}

	logInf.Println("JSON: "+artJson.IsNew+", "+artJson.Title)
	if "yes" == artJson.IsNew {
		logInf.Println("Inserting into db a new article")
		//处理新文章入库的流程
		model.AddArticleByInf(artJson.Title, artJson.Content,1)
		//todo:入库失败校验以及对应的状态码返回
	}else{
		logInf.Println("Updating an existing article in db: ")
		//todo:处理编辑旧文章的页面
	}

	//上传文章的ajax请求要的内容是json，所以这里也应该返回json而不是html
	//要符合格式：
	c.JSON(http.StatusOK,gin.H{"status":0, "message":"Posted an article", "data":nil})

	logInf.Println("Leaving EditArticlePost for article: ",artJson.Title)
}

//readArticle的GET请求（获取对应id的文章内容）
func ReadArticleGet(c *gin.Context){
	logInf.Println("Entering ReadArticleGet: ")

	//先获取文章id并转int64
	idStr:=c.Request.FormValue("id")
	id,_ := strconv.ParseInt(idStr,10,64)

	//通过ID查表
	art := model.Article{ID:id}
	model.GetArticleById(id, &art)

	//检验结果
	if art.Title == "" {
		logErr.Println("cannot find article by id: ", id)
	}

	c.HTML(http.StatusOK,"readArticle.html",gin.H{
		"Id": art.ID,
		"Title": art.Title,
		"Content": art.Content,
	})

	logInf.Println("Leaving ReadArticleGet for article: ", art.Title, " with content: ", art.Content)
}

//deleteArticle的POST请求（删除指定id的文章）
func DeleteArticlePost(c *gin.Context){
	logInf.Println("Entering DeleteArticlePost: ")

	//这里是从url解析参数而不是解析json，别搞错了
	//先获取文章id并转为int64
	idStr:=c.Request.FormValue("id")
	id,_ := strconv.ParseInt(idStr,10,64)

	//交给model模块删除之
	logInf.Println("Deleting article with id:", idStr)
	model.DeleteArticle(id)

	//检查model模块的全局变量Error，看看是否有异常发生
	if model.Error{
		//返回json的统一格式：{code:200,obj:{status:0(成功)/1（失败）,message:"content",data:"content"}}
		c.JSON(http.StatusOK,gin.H{"status":1,"message":"Error deleting article","data":""})
	}else{
		c.JSON(http.StatusOK,gin.H{"status":0,"message":"Article deleted","data":""})
	}


	logInf.Println("Leaving DeleteArticlePost for article id = : ", id)
}