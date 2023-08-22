package controllers

import (
	"ginapi/mysqlexample"
	"ginapi/structs"
	"ginapi/utils/token"
	"net/http"

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
