package router

import (
	"ginapi/controllers"
	"ginapi/esexample"
	"ginapi/middlewares"

	"github.com/gin-gonic/gin"
)

func PublicRouter(c *gin.Engine) {
	public := c.Group("/api")
	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)
	public.GET("/baidu", controllers.RedirectLink)
}

func ProtectRouter(c *gin.Engine) {
	protected := c.Group("/admin")
	//中间件
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", controllers.CurrentUser)
	protected.POST("/upload", controllers.UploadFile)
	protected.POST("/bodytostruct", controllers.BodyToStruct)
	protected.POST("/bodytomap", controllers.BodyToMap)
}

func EsRouter(c *gin.Engine) {
	public := c.Group("/es")
	public.POST("/get", esexample.DynamicDSL)
}
