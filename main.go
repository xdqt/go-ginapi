package main

import (
	"ginapi/mysqlexample"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"ginapi/controllers"

	"ginapi/esexample"
	"ginapi/ossexample"
	"ginapi/router"
)

func main() {

	file, _ := os.OpenFile("ellis.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
	log.SetOutput(file)
	defer file.Close()
	// optional: log date-time, filename, and line number
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//数据库初始化以及建表
	mysqlexample.Initdb()
	ossexample.InitS3()
	esexample.InitEs()
	// esexample.SearchByDSL()
	// // esexample.IndexOneDocument()
	// esexample.UpdateByQuery()
	// esexample.Delete()
	// esexample.SearchAfter()
	esexample.NestedQueryVersion2()

	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.PublicRouter(r)
	router.ProtectRouter(r)
	router.EsRouter(r)

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
