package controllers

import (
	"fmt"
	"ginapi/mysqlexample"
	"ginapi/structs"
	"ginapi/utils/token"
	"net/http"
	"reflect"

	"ginapi/ossexample"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	defer f.Close()
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

func RedirectLink(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "https://www.baidu.com")
}

func GetQuery(c *gin.Context) {
	s, ok := c.GetQuery("name")
	if ok {
		fmt.Printf("s: %v\n", s)
	} else {
		fmt.Println("没有传递这个参数")
	}

	values, _ := c.GetQueryArray("arrayfield")

	fmt.Printf("values: %v\n", values)

	c.JSON(http.StatusOK, gin.H{"name": s})
}

func DynamicParams(c *gin.Context) {
	user := c.Param("user_id")
	fmt.Printf("user: %v\n", user)
	book := c.Param("book_id")
	fmt.Printf("book: %v\n", book)
	c.JSON(http.StatusOK, gin.H{"user": user, "book": book})
}

func Validation(c *gin.Context) {
	var user structs.UserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func CustomError(err error, obj any) map[string]string {
	var errors map[string]string = make(map[string]string)

	getObj := reflect.TypeOf(obj)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, v := range errs {
			if sf, exist := getObj.Elem().FieldByName(v.Field()); exist {
				errors[v.Field()] = sf.Tag.Get("msg")
			}
		}
		return errors
	} else {
		return nil
	}

}

func Validationcustomerror(c *gin.Context) {
	var user structs.UserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		errors := CustomError(err, &user)
		if errors != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		} else {
			c.JSON(http.StatusOK, gin.H{"MSG": user})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"MSG": user})
	}
}
