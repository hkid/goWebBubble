package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)
//Todo model
type Todo struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Status bool `json:"status"`
}
var (
	DB *gorm.DB
)
func initMysql() (err error){
	dsn:="root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	DB,err =gorm.Open(mysql.Open(dsn),&gorm.Config{}) //这不要 :=
	if err!=nil{
		return err
	}
	return nil
}
func main() {
	// 0 创建数据库
	// 1 连接数据库 ->封装为函数
	err:=initMysql()
	if err!=nil{
		panic(err)
	}
	//模型绑定
	DB.AutoMigrate(&Todo{})


	router:=gin.Default()
	//router.LoadHTMLGlob("./templates/*")
	router.LoadHTMLFiles("./templates/index.html")   //解析模板
	router.Static("/static","./static")	   //加载静态文件
	router.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK,"index.html",nil) //渲染模板
	})
	//待办事项 v1
	v1Group:=router.Group("v1")
	{  //这个括号必须下来
	//添加
	v1Group.POST("/todo", func(context *gin.Context) {
		//1 从请求中把数据拿出来  BindJSON
		//2 存入数据库
		//3 返回响应
		var todo Todo
		context.BindJSON(&todo)

		/*if err=DB.Create(&todo).Error;err!=nil{
			context.JSON(...)
		}*/

		err:=DB.Create(&todo).Error
		if err!=nil{
			context.JSON(http.StatusOK,gin.H{"error":err.Error()})
		}else {
		/*	context.JSON(http.StatusOK,gin.H{
				"code":2000,
				"msg":"success",
				"data":todo,
			}}*/
			context.JSON(http.StatusOK,todo)
		}
	})
	//删除
	v1Group.DELETE("/todo/:id", func(context *gin.Context) {
	/*
		/user/search/小王子/
		"/user/search/:name"
		name:=context.Param("name")
		*/
		id:=context.Param("id")

		if err =DB.Where("id=?",id).Delete(Todo{}).Error;err!=nil{
			context.JSON(http.StatusOK,gin.H{"error":err.Error()})
		}else {
			context.JSON(http.StatusOK,gin.H{"success":"deleted"})
		}
	})
	//改
	v1Group.PUT("/todo/:id", func(context *gin.Context) {
		id:=context.Param("id")
		var todo Todo
		if err=DB.Where("id=?",id).First(&todo).Error;err!=nil{
			context.JSON(http.StatusOK,gin.H{
				"error":err.Error(),
			})
		}else{
			context.BindJSON(&todo)
			if err=DB.Save(&todo).Error;err!=nil{
				context.JSON(http.StatusOK,gin.H{
					"error":err.Error(),
				})
			}else {
				context.JSON(http.StatusOK,todo)
			}
		}
	})
	//查看所有的待办事项
	v1Group.GET("/todo", func(context *gin.Context) {
		//查看表中所有的数据
		var todoList []Todo
		 if err=DB.Find(&todoList).Error;err!=nil{//失败
		 	context.JSON(http.StatusOK,gin.H{"error":err.Error()})
		 }else {  //成功
		 	context.JSON(http.StatusOK,todoList)
		 }
	})
	//查看某个待办事项
	v1Group.GET("/todo/:id", func(context *gin.Context) {
		
	})

	}
	router.Run(":80")
}
