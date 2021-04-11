package main

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"trueNewblog/sql"
)

func main() {
	r := gin.Default()
	sql.InitDB()//账号数据库
	r.POST("/register",register)//注册
	r.POST("/login",login)//登录
	r.POST("/postblogs",upload)
	//r.POST("/register",postblogs)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
func register(c *gin.Context ){
	flag :=sql.InsertRowDemo(c)
	switch flag{
	case -1:{
		c.String(501,fmt.Sprintf("注册失败！密码不一致.\n"))
	}
	case 0:{
		c.String(501,fmt.Sprintf("注册失败！用户名重复！\n"))
	}
	case -2:{
		c.String(501,fmt.Sprintf("注册失败！内部错误！\n"))
	}
	default:{
		c.String(http.StatusOK,fmt.Sprintf("注册成功, 您是第 %d个用户.\n",flag))
	}
	}
	}
func login(c *gin.Context ){
	flag :=sql.QueryRowDemo(c)
	if flag ==1{
		c.String(http.StatusOK,fmt.Sprintf("登录成功！\n"))
		cookie, err := c.Cookie("gin_cookie") // 获取Cookie
		if err != nil {
			cookie = "NotSet"
			// 设置Cookie
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}
		fmt.Printf("Cookie value: %s \n", cookie)
	} else{
		c.String(501,fmt.Sprintf("登录失败!\n"))
	}
}
func upload(c *gin.Context){
	//创建一个新文件，写入内容
	filePath := "./output.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("打开文件错误= %v \n", err)
		return
	}
	//及时关闭
	defer file.Close()
	//写入内容
	str := c.PostForm("text") // \n\r表示换行  txt文件要看到换行效果要用 \r\n
	//写入时，使用带缓存的 *Writer
	writer := bufio.NewWriter(file)
	writer.WriteString(str)
	//因为 writer 是带缓存的，因此在调用 WriterString 方法时，内容是先写入缓存的
	//所以要调用 flush方法，将缓存的数据真正写入到文件中。
	writer.Flush()
	//file, _ :=c.FormFile("file")//取文件
	c.String(http.StatusOK,fmt.Sprintf("finished uploading! "))
}
//func postblogs(c *gin.Context ){
//	test := sql.QueryRowDemo(c)
//	if test ==1{
//		c.String(http.StatusOK,fmt.Sprintf("登录成功！\n"))
//	}else{
//		c.String(501,fmt.Sprintf("登录失败!\n"))
//	}
//}
