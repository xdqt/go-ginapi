package controllers

import (
	"ginapi/mysqlexample"
	"ginapi/structs"
	"ginapi/utils/token"
	"net/http"

	"ginapi/ossexample"

	"github.com/gin-gonic/gin"
)

// func PublicRoutes(route *gin.RouterGroup) {
// 	route.POST("/register", Register)
// 	route.POST("/login", Login)
// }

// func ProtectedRoutes(route *gin.RouterGroup) {
// 	route.POST("/user", CurrentUser)
// }

func Register(c *gin.Context) {

	var input structs.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := structs.User{UserName: input.Username, Email: "email"}

	result := mysqlexample.CheckUserExist(input.Username, &user)
	if result {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户已经存在"})
		return
	} else {
		mysqlexample.Create(&user)
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success!"})
}

func Login(c *gin.Context) {

	var input structs.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := structs.User{}

	u.UserName = input.Username

	result, token := mysqlexample.LoginCheck(u.UserName)

	if !result {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func CurrentUser(c *gin.Context) {
	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user structs.User
	mysqlexample.GetUserByID(int(user_id), &user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": user})
}

func UploadFile(c *gin.Context) {

	file, err := c.FormFile("file")
	bucketname := c.PostForm("bucketname")
	// c.SaveUploadedFile(file, file.Filename)
	// key := c.PostForm("key")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	f, err2 := file.Open()
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	ossexample.UploadFileObj(f, bucketname, file.Filename)
	f.Close()
	c.JSON(http.StatusOK, gin.H{"meta": "success"})
}

func BodyToStruct(c *gin.Context) {
	var json structs.EsQuery
	if err := c.ShouldBind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"value": json})
	}
}

func BodyToMap(c *gin.Context) {
	// var json structs.EsQuery
	var json map[string]interface{}
	if err := c.ShouldBind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"value": json})
	}
}
