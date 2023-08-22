package main

import (
	"ginapi/mysqlexample"

	"github.com/gin-gonic/gin"

	"ginapi/controllers"

	"ginapi/middlewares"
)

func main() {
	//数据库初始化以及建表
	mysqlexample.Initdb()
	r := gin.Default()

	public := r.Group("/api")
	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)

	protected := r.Group("/admin")
	//中间件
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", controllers.CurrentUser)
	for _, item := range r.Routes() {
		println("method:", item.Method, "path:", item.Path)
	}
	r.Run(":8080")

}
