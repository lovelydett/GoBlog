package model

import(
	//"github.com/jinzhu/gorm"
	"time"
)


type Category struct {
	Id int64 `orm:"pk"`
	Name string `json:"name"`
	Articles []*Article `orm:"reverse(many);null"` //分类允许没有文章
	//NumofArticles int64 ?
}

//文章id要使用uuid!否则删除不方便，然后在id上建索引
func getTimeStamp () int64 {
	time.Sleep(2*time.Millisecond)//防止重复id，需要吗？看uuid原理
	return time.Now().Unix()
}
type Article struct {
	Id int64 `orm:"pk;index"`
	Title string `json:"title"`
	Content string `json:"content"`
	// TimeInit int64 `orm:"index"` 不需要 uuid就是生成时间戳
	TimeLastEdit int64 `orm:"index"`
	ReaderCount int64 `orm:"index"`
	Categories []*Category `orm:"index;rel(m2m);null"` //指针数组，文章可以没有标签
}