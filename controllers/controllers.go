package controllers

import (
	"GinBlog/model"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 常量
const (
	PAGE_SIZE = 10
)

// 日志系统
var logInf = log.New(os.Stdout, "[INFO]", log.LstdFlags)
var logErr = log.New(os.Stdout, "[Error]", log.LstdFlags|log.Lshortfile)

// 上取整
func ceil(n1, n2 int64) int64 {
	return int64(math.Ceil(float64(n1) / float64(n2)))
}

// cookie判断
func isLogin(c *gin.Context) bool {
	_, e := c.Cookie("adminCookie")
	return e == nil
}

// int64 to int
func int64ToInt(x int64) int {
	strInt64 := strconv.FormatInt(x, 10)
	ret := 0
	if r, e := strconv.Atoi(strInt64); e == nil {
		ret = r
	}
	return ret
}

// home页面的GET请求
func HomeGet(c *gin.Context) {
	logInf.Println("Entering HomeGet")

	// 打包分類標籤
	categories := []model.Category{}
	model.GetAllCategories(&categories)

	// 打包文章, 首頁顯示最新的10篇文章
	articles := []model.Article{}
	model.GetPagedArticles(0, PAGE_SIZE, &articles)

	c.HTML(http.StatusOK, "home.html", gin.H{
		"Categories": categories,
		"Articles":   articles,
	})

	logInf.Println("Leaving HomeGet")
}

//login页面的GET
func LoginGet(c *gin.Context) {
	logInf.Println("Entering LoginGet")

	c.HTML(http.StatusOK, "login.html", nil)
	//c.String(http.StatusOK,"ok")

	logInf.Println("Leaving LoginGet")
}

//404页面的GET请求
func NotFoundGet(c *gin.Context) {
	logInf.Println("Entering NotFoundGet for: " + c.Request.URL.String())

	c.HTML(http.StatusNotFound, "404.html", nil)

	logInf.Println("Leaving NotFoundGet")
}

// article页面的GET请求
func ArticleGet(c *gin.Context) {
	logInf.Println("Entering ArticleGet")

	// 獲取當前文章總數
	numArticles := model.GetNumArticles()
	numPage := ceil(numArticles, PAGE_SIZE)

	// 解析url參數pageNum
	pageNum_s := c.Request.FormValue("pageNum")
	pageNum := int64(0)
	if len(pageNum_s) > 0 {
		// 有這個參數且成功解析
		if pn, err := strconv.ParseInt(pageNum_s, 10, 64); err == nil {
			pageNum = pn - 1 // 前端[1, numPage] 後端[0, numPage-1]
		}
		// 檢查文章總頁數, 如果pageNum溢出則設置爲0
		if pageNum < 0 || pageNum >= numPage {
			pageNum = 0
		}
	}

	// 獲取第pageNum頁的文章
	articles := []model.Article{}
	model.GetPagedArticles(pageNum, PAGE_SIZE, &articles)

	// 獲取分类列表
	categories := []model.Category{}
	model.GetAllCategories(&categories)

	// handle登录状态
	login := isLogin(c)

	// 打包模板内容
	c.HTML(http.StatusOK, "article.html", gin.H{
		"IsLogin":     login,
		"Categories":  categories,
		"Articles":    articles,
		"NumArticles": int64ToInt(numArticles),
		"NumPage":     int64ToInt(numPage),
		"PageNum":     int64ToInt(pageNum) + 1,
	})

	logInf.Println("Leaving ArticleGet for ", pageNum, "/", numPage)
}

//editArticle页面的GET请求
func EditArticleGet(c *gin.Context) {
	logInf.Println("Entering EditArticleGet for article: ")
	isNew := ("yes" == c.Request.FormValue("isNew"))
	article := &model.Article{}
	if !isNew {
		//处理编辑旧文章的页面
		logInf.Println("editing an existing article: ")

		//解析url参数，先获取文章id并转为int64
		idStr := c.Request.FormValue("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		//查表，获取文章title和content
		article.ID = id
		model.GetArticleById(id, article)
	}

	//模板传参
	c.HTML(http.StatusOK, "editArticle.html", gin.H{
		"IsNew":   isNew,
		"Id":      article.ID,
		"Title":   article.Title,
		"Content": article.Content,
	})

	logInf.Println("Leaving EditArticleGet for article: ")
}

type EditArtJson struct {
	IsNew   string `json:"isNew"`
	IdStr   string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

//editArticle页面的POST请求
func EditArticlePost(c *gin.Context) {
	logInf.Println("Entering EditArticlePost: ")

	//需要解析JSON
	artJson := EditArtJson{}
	if nil != c.BindJSON(&artJson) {
		logErr.Println("unable to unmarshal JSON posted from editArticle.html")
		logInf.Println("Leaving EditArticlePost for article: ")
		c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Unable to post an article, bind JSON error", "data": nil})
		return
	}

	logInf.Println("JSON from editArticle: "+artJson.IdStr, ",", artJson.IsNew, ",", artJson.Title)
	if "yes" == artJson.IsNew {
		logInf.Println("Inserting into db a new article")
		//处理新文章入库的流程
		model.AddArticleByInf(artJson.Title, artJson.Content)
		//todo:入库失败校验以及对应的状态码返回
	} else {
		logInf.Println("Updating an existing article in db:", artJson.Title)
		//处理编辑旧文章的页面
		id, _ := strconv.ParseInt(artJson.IdStr, 10, 64)
		article := &model.Article{ID: id, Title: artJson.Title, Content: artJson.Content}
		model.UpdateArticle(id, article)
		if model.Error {
			c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Unable to update article " + article.Title, "data": nil})
			return
		}
	}

	//上传文章的ajax请求要的内容是json，所以这里也应该返回json而不是html
	//要符合格式：
	c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Posted an article", "data": nil})
	logInf.Println("Leaving EditArticlePost for article: ", artJson.Title)
}

//readArticle的GET请求（获取对应id的文章内容）
func ReadArticleGet(c *gin.Context) {
	logInf.Println("Entering ReadArticleGet: ")

	//先获取文章id并转int64
	idStr := c.Request.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	//通过ID查表
	art := model.Article{ID: id}
	model.GetArticleById(id, &art)

	//检验结果
	if art.Title == "" {
		logErr.Println("cannot find article by id: ", id)
	}

	c.HTML(http.StatusOK, "readArticle.html", gin.H{
		"IsLogin": isLogin(c),
		"Id":      art.ID,
		"Title":   art.Title,
		"Content": art.Content,
	})

	logInf.Println("Leaving ReadArticleGet for article: ", art.Title, " with content: ", art.Content)
}

//deleteArticle的POST请求（删除指定id的文章）
func DeleteArticlePost(c *gin.Context) {
	logInf.Println("Entering DeleteArticlePost: ")

	//这里是从url解析参数而不是解析json，别搞错了
	//先获取文章id并转为int64
	idStr := c.Request.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	//交给model模块删除之
	logInf.Println("Deleting article with id:", idStr)
	model.DeleteArticle(id)

	//检查model模块的全局变量Error，看看是否有异常发生
	if model.Error {
		//返回json的统一格式：{code:200,obj:{status:0(成功)/-1（失败）,message:"content",data:"content"}}
		c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Error deleting article", "data": ""})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Article deleted", "data": ""})
	}

	logInf.Println("Leaving DeleteArticlePost for article id = : ", id)
}

//login的POST请求（登录）
func LoginPost(c *gin.Context) {
	logInf.Println("Entering LoginPost")

	//从url请求中获取密码
	pwStr := c.Request.FormValue("pw")
	if strings.Compare(pwStr, "981123") == 0 {
		//设置一条名叫adminCookie的cookie
		c.SetCookie("adminCookie", "tt", 3600, "/", "localhost", false, true)
		c.JSON(200, gin.H{"status": 0, "message": "Login Successful", "data": ""})
	} else {
		c.JSON(200, gin.H{"status": 1, "message": "Wrong Password", "data": ""})
	}

	logInf.Println("Leaving LoginPost")
}
