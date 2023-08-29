package structs

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

// register 路由body验证结构体
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录body数据验证结构体
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	gorm.Model
	UserName string `gorm:"column:name;size:255;unique"`
	Email    string `gorm:"column:email"`
}

type Tabler interface {
	TableName() string
}

// 设置表名称
func (recv User) TableName() string {
	return "ellis_user"
}

type YamlStruct struct {
	Mysql     Mysql  `yaml:"Mysql"`
	JWTSecret string `yaml:"JWTSecret"`
}

type Mysql struct {
	Host      string `yaml:"Host"`
	Port      int    `yaml:"Port"`
	UserName  string `yaml:"UserName"`
	Password  string `yaml:"Password"`
	Database  string `yaml:"Database"`
	Charset   string `yaml:"Charset"`
	ParseTime bool   `yaml:"ParseTime"`
	Loc       string `yaml:"Loc"`
}

type EsQuery struct {
	Country   string `json:"country"`
	EmailAddr string `json:"emailAddr"`
}

type MongoStruct struct {
	Id       primitive.ObjectID `bson:"_id"`
	UserName string             `bson:"username"`
	Email    string             `bson:"email"`
}

type UserInfo struct {
	UserName  string `json:"username" binding:"required,min=4,max=6" msg:"username验证失败"`
	Age       int    `json:"age" binding:"gt=18,lte=120" msg:"age验证失败"`
	Password  string `json:"password" binding:"required" msg:"密码验证失败"`
	Password2 string `json:"password2" binding:"required,eqfield=Password" msg:"必填且需要与password一致"`
}
