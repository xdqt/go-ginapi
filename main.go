package main

import (
	"ginapi/mysqlexample"

	"github.com/gin-gonic/gin"

	"ginapi/controllers"

	"ginapi/middlewares"

	"ginapi/ossexample"
)

func main() {
	//数据库初始化以及建表
	mysqlexample.Initdb()
	ossexample.InitS3()
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	public := r.Group("/api")
	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)
	public.GET("/baidu", controllers.RedirectLink)

	protected := r.Group("/admin")
	//中间件
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", controllers.CurrentUser)
	protected.POST("/upload", controllers.UploadFile)
	protected.POST("/bodytostruct", controllers.BodyToStruct)
	protected.POST("/bodytomap", controllers.BodyToMap)

	testquery := r.Group("/query")
	// 查询参数
	testquery.GET("getquery", controllers.GetQuery)
	// http://localhost:8080/query/getquery?name=ellis&arrayfield=1&arrayfield=2
	// 动态参数
	testquery.GET("dynamicparams/:user_id/:book_id", controllers.DynamicParams)
	// http://localhost:8080/query/dynamicparams/1/1
	testquery.POST("/validation", controllers.Validation)

	testquery.POST("/validationcustomerror", controllers.Validationcustomerror)
	for _, item := range r.Routes() {
		println("method:", item.Method, "path:", item.Path)
	}

	r.Run(":8080")
}
