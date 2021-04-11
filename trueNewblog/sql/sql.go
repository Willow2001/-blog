package sql

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)
// 定义一个全局对象db
var db *sql.DB
//type article struct{
//	id   int
//	name string
//	article string
//	author string
//	LikeNum int
//}
type user struct {
	id   int
	age  string//密码
	name string
	artiid int
}
////点赞
//func (a *Author) Like(v *VideoInformation) {
//	v.LikeNum++
//	fmt.Println(a.Name,"点赞了该视频")
//}
//发表文章
//func Insertblog(c *gin.Context ) int {
//	sqlStr1 := "insert into article(name, article) values (?,?)"//插入语句
//	n := c.PostForm("username")
//	k1 := c.PostForm("password1")
//	k2 := c.PostForm("password2")
//	if k1!=k2{
//		return -1
//	}
//	flag := QueryRowDemo(c)
//	if flag != -2{//找到了该元组
//		return 0
//	}
//	ret, err := db.Exec(sqlStr1, n, k1)
//	if err != nil {
//		fmt.Printf("insert failed, err:%v\n", err)
//		return -2
//	}
//	theID, err := ret.LastInsertId() // 新插入数据的id
//	if err != nil {
//		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
//		return -2
//	}
//	return int(theID)
//}
// 插入数据(注册）
func InsertRowDemo(c *gin.Context ) int {
	sqlStr1 := "insert into user(name, age) values (?,?)"//插入语句
	n := c.PostForm("username")
	k1 := c.PostForm("password1")
	k2 := c.PostForm("password2")
	if k1!=k2{
		return -1
	}
	flag := QueryRowDemo(c)
	if flag != -2{//找到了该元组
		return 0
	}
	ret, err := db.Exec(sqlStr1, n, k1)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return -2
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return -2
	}
	return int(theID)
}
// 查询账号密码是否匹配
func QueryRowDemo(c *gin.Context) int{
	var testflag = 0
	sqlStr := "select name, age from user where name=?"
	var u user
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	n := c.PostForm("username")
	k := c.PostForm("password1")
	fmt.Printf(c.PostForm("username"))
	err := db.QueryRow(sqlStr,n).Scan(&u.name, &u.age)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return -2//未找到
	}
	if k == u.age{
		testflag = 1//账号密码匹配
	}
	fmt.Printf("name:%s age:%s\n", u.name, u.age)
	return testflag
}
// 定义一个初始化数据库的函数
func InitDB() (err error) {
	// DSN:Data Source Name
	dsn := "root:151413Mn++@tcp(127.0.0.1:3306)/sql_blog?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil {
		return err
	} else{
		fmt.Printf("成功连接\n")
	}
	return nil
}

func main() {
	//var testflag,n int
	//err := initDB() // 调用输出化数据库的函数
	//if err != nil {
	//	fmt.Printf("init db failed,err:%v\n", err)
	//	return
	//}
	//fmt.Printf("请选择功能\n")
	//fmt.Printf("1.登录\n")
	//fmt.Printf("2.注册\n")
	//fmt.Scan(&n)
	//switch n{
	//case 1: {
	//	testflag = queryRowDemo()
	//	for testflag != 1 {
	//		fmt.Printf("账号或密码错误，请重试\n")
	//		testflag = queryRowDemo()
	//	}
	//	fmt.Printf("账号密码正确,成功登录\n")
	//}
	//case 2:{
	//	insertRowDemo()
	//}
	//}
}