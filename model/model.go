package model

import (
	. "GinBlog/utils"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	// _ "github.com/jinzhu/gorm/dialects/mysql"
)

//错误标志位，类似gorm的方式来指示操作有无成功
var Error bool

//日志
var logInf = log.New(os.Stdout, "[INFO]", log.LstdFlags)
var logErr = log.New(os.Stdout, "[Error]", log.LstdFlags|log.Lshortfile)

// sql表名: categories
type Category struct {
	ID   int64  `gorm:"primary_key"`
	Name string `json:"name"`
}

//文章id要使用uuid!否则删除不方便，然后在id上建索引
func getTimeStamp() int64 {
	time.Sleep(2 * time.Millisecond) //防止重复id，需要吗？
	return time.Now().Unix()
}

// sql表名: articles
type Article struct {
	ID      int64  `gorm:"index;primary_key;"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Preview string `json:"preview"`
	// TimeInit int64 `orm:"index"` 不需要 uuid就是生成时间戳
	TimeLastEdit int64      `gorm:"index;"`
	ReaderCount  int64      `gorm:"index;"`
	Categories   []Category `gorm:"index;many2many:article_category;"` //文章包含并属于多个分类
}

//初始化model模块：数据库连接，校验等
var db *gorm.DB

func init() {
	logInf.Println("Enter init of model.go")

	//错误标志置为否
	Error = false

	//打开链接
	dbType := "sqlite3"
	dbName := "database/blogDB.db"
	db_, err := gorm.Open(dbType, dbName)

	//校验数据库连接
	if err != nil {
		panic("Enable to open DB: " + dbType + ", " + dbName)
	}

	//赋值给全局变量
	db = db_

	//迁移
	db.AutoMigrate(&Category{}, &Article{})

	//插入默认数据
	cat := Category{ID: 1, Name: "default"}
	if db.NewRecord(cat) {
		//还没有分类
		logInf.Println("Auto create default category")
		db.Create(&cat)
	}

	article := Article{ID: 1, Title: "Default title", Content: "Default content", Preview: "Default preview"}
	if db.NewRecord(article) {
		//还没有文章
		logInf.Println("Auto create default article")
		addArticle(&article)
	}

	logInf.Println("Leave init of model.go")
}

func closeDB() {
	if nil != db.Close() {
		logErr.Println("Enable to close db")
	}
}

//CRUD方法

// 统计数据查询
// 查询文章总数
func GetNumArticles() int64 {
	var ret int64 = 0
	db.Table("articles").Count(&ret)
	logInf.Println("Ariticle counted: ", ret)
	return ret
}

// 查询分类总数
func GetNumCategories() int64 {
	var ret int64 = 0
	db.Table("categories").Count(&ret)
	logInf.Println("Category counted: ", ret)
	return ret
}

//查询
//通过ID查文章
func GetArticleById(id int64, article *Article) {
	db.First(article, id)
}

//通过精确名称查找文章
func GetArticleByTitle(title string, article *Article) {
	db.Where("title = ?", title).First(article)
}

// 查找所有文章
func GetAllArticles(articles *[]Article) {
	logInf.Println("Enter GetAllArticles")
	db.Find(articles)
	logInf.Println("Leave GetAllArticles")
}

// 查找指定頁面的文章
func GetPagedArticles(pageNum, pageSize int64, articles *[]Article) {
	logInf.Println("Enter GetPagedArticles at ", pageNum)
	if pageNum < 0 || pageSize < 1 {
		logErr.Println("Invalid pageNum,pageSize: ", pageNum, pageSize)
	}
	db.Limit(pageSize).Offset(pageNum * pageSize).Order("id desc").Find(articles)
	logInf.Println("Leave GetPagedArticles")
}

// 查找指定分类
func GetCategoryByName(name string, category *Category) {
	logInf.Println("Enter GetCategoryByName")
	db.Where("name = ?", name).First(category)
	logInf.Println("Leave GetCategoryByName")
}

// 查找所有分类
func GetAllCategories(categories *[]Category) {
	logInf.Println("Enter GetAllCategories")
	db.Find(categories)
	logInf.Println("Leave GetAllCategories")
}

// 按分類查找文章
func GetArticlesByCategory(category *Category, articles *[]Article) {
	logInf.Println("Enter GetArticleByCategory")
	db.Model(category).Association("Article").Find(articles)
	logInf.Println("Leave GetArticleByCategory")
}

// 查找一篇文章的所有分類
func GetCategoriesOfArticle(article *Article, categories *[]Category) {
	logInf.Println("Enter GetCategoriesOfArticle")
	// Todo: in model.go the column names "Categories", however it is mapped to "category" in database
	if err := db.Model(article).Association("Categories").Find(categories).Error; err != nil {
		logErr.Println("Unable to get categories for article:", article.Title, ",", err.Error())
	}
	logInf.Println("Leave GetCategoriesOfArticle")
}

// 增加
// 增加一篇文章
func AddArticleByInf(title string, content string) {
	logInf.Println("Enter AddArticleByInf: " + title)
	// 文章默認沒有任何分類
	art := Article{
		ID:      getTimeStamp(),
		Title:   title,
		Content: content,
		Preview: content[:Min(len(content), PREVIEW_STRING_LENGTH)],
	}
	addArticle(&art)
	logInf.Println("Leave AddArticleByInf: " + title)
}

func addArticle(article *Article) {
	db.Create(article)
}

// 增加一個分類
func AddCategoryByName(name string) {
	logInf.Println("Enter AddCategoryByName: " + name)
	db.Create(&Category{Name: name})
	logInf.Println("Leave AddCategoryByName: " + name)
}

// 增加一條文章-分類關系
func AddArticleCategoryAssociation(article *Article, category *Category) {
	logInf.Println("Enter AddArticleCategoryAssociation")
	db.Model(article).Association("Categories").Append(category)
	logInf.Println("Leave AddArticleCategoryAssociation")
}

// 删除
// 通过id删除文章
func DeleteArticle(id int64) {
	logInf.Println("Enter DeleteArticle")

	article := &Article{ID: id}
	//检查是否还有这个id的文章
	//todo:注意：这是gorm检查错误情况的标准写法，需要都改成这样
	if err := db.Delete(article).Error; err != nil {
		logInf.Println("Unable to delete article id = ", id)
		Error = true
		return
	}

	Error = false
	logInf.Println("Leave DeleteArticle")
}

// Delete a category
func DeleteCategory(name string) {
	logInf.Println("Deleting category: ", name)
	if err := db.Delete(Category{}, "name = ?", name).Error; err != nil {
		logInf.Println("Unable to delete category: ", name)
		Error = true
		return
	}
	Error = false
}

//更新
//通过id更新文章
func UpdateArticle(id int64, article *Article) {
	logInf.Println("Enter UpdateArticle for:", article.Title)
	if err := db.Save(article).Error; err != nil { // Save update value in database, if the value doesn't have primary key, will insert it
		logInf.Println("Enable to update article:", article.Title)
		Error = true
		return
	}
	Error = false
	logInf.Println("Leave UpdateArticle")
}
