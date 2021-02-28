package model

import(
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"

	// _ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

//错误标志位，类似gorm的方式来指示操作有无成功
var Error bool

//日志
var logInf = log.New(os.Stdout, "[INFO]", log.LstdFlags)
var logErr = log.New(os.Stdout, "[Error]", log.LstdFlags | log.Lshortfile)


type Category struct {
	ID int64 `gorm:"primary_key"`
	Name string `json:"name"`
	//NumofArticles int64 ?
}

//文章id要使用uuid!否则删除不方便，然后在id上建索引
func getTimeStamp () int64 {
	time.Sleep(2*time.Millisecond)//防止重复id，需要吗？看uuid原理
	return time.Now().Unix()
}
type Article struct {
	ID int64 `gorm:"index;primary_key;"`
	Title string `json:"title"`
	Content string `json:"content"`
	// TimeInit int64 `orm:"index"` 不需要 uuid就是生成时间戳
	TimeLastEdit int64 `gorm:"index;"`
	ReaderCount int64 `gorm:"index;"`
	Categories []Category `gorm:"index;many2many:article_category;"` //文章包含并属于多个分类
}

//初始化model模块：数据库连接，校验等
var db *gorm.DB
func init(){
	logInf.Println("Enter init of model.go")

	//错误标志置为否
	Error = false

	//打开链接
	dbType := "sqlite3"
	dbName := "database/blogDB.db"
	db_, err := gorm.Open(dbType, dbName)

	//校验数据库连接
	if err!=nil {
		//logErr.Println("Enable to open DB")
		panic("Enable to open DB: "+dbType+", "+dbName)
	}

	//赋值给全局变量
	db = db_

	//迁移
	db.AutoMigrate(&Category{},&Article{})

	//插入默认数据
	cat := Category{ID:1, Name:"default"}
	if db.NewRecord(cat) {
		//还没有分类
		logInf.Println("Auto create default category")
		db.Create(&cat)
	}

	article := Article{ID:1, Title: "Default title", Content: "Default content"}
	if db.NewRecord(article) {
		//还没有文章
		logInf.Println("Auto create default article")
		AddArticle(&article)
	}

	logInf.Println("Leave init of model.go")
}

func closeDB(){
	if nil!=db.Close(){
		logErr.Println("Enable to close db")
	}
}

//CRUD方法

//查询
//通过ID查文章
func GetArticleById(id int64, article *Article){
	db.First(article, id)
}

//通过精确名称查找文章
func GetArticleByTitle(title string, article *Article){
	db.Where("title = ?", title).First(article)
}

func GetArticles(num int, articles *[]Article){
	logInf.Println("Enter GetArticle")
	//if num<=0, get all
	db.Find(articles)
	//select * from article

	logInf.Println("Leave GetArticle")
}


func AddArticleByInf(title string, content string, catId int64){
	logInf.Println("Enter AddArticleByInf: "+title)
	//cats := []Category{{ID:catId}}
	art := Article{ID:getTimeStamp(),Title: title, Content: content}
	AddArticle(&art)
	logInf.Println("Leave AddArticleByInf: "+title)
}

func AddArticle(article *Article) {
	logInf.Println("Enter AddArticle")
	db.Create(article)
	logInf.Println("Leave AddArticle")
}

//删除
//通过id删除文章
func DeleteArticle(id int64){
	logInf.Println("Enter DeleteArticle")

	article := &Article{ID:id}
	//检查是否还有这个id的文章
	//todo:注意：这是gorm检查错误情况的标准写法，需要都改成这样
	if err:=db.Delete(article).Error; err!=nil {
		logInf.Println("Enable to delete article id = ", id)
		Error = true
		return
	}

	Error = false
	logInf.Println("Leave DeleteArticle")
}




