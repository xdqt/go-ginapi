package structs

import (
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
