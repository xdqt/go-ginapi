package mysqlexample

import (
	"fmt"
	"ginapi/structs"
	tokentool "ginapi/utils/token"
	"io/ioutil"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

func readConf(filename string) (*structs.YamlStruct, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := structs.YamlStruct{}
	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}
	fmt.Printf("c: %v\n", c)
	fmt.Printf("c.Mysql.Host: %v\n", c.Mysql.Host)
	return &c, err
}

func Initdb() {
	// str, _ := os.Getwd()
	// fmt.Printf("str: %v\n", str)

	// _path := filepath.Join(str, "config", "config.yaml")
	// fmt.Printf("_path: %v\n", _path)

	c, _ := readConf("./config/config.yaml")

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=%v&loc=Local", c.Mysql.UserName, c.Mysql.Password, c.Mysql.Host, c.Mysql.Port, c.Mysql.Database, c.Mysql.Charset, c.Mysql.ParseTime)

	// dsn := "ellis:ellis@tcp(192.168.214.134:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	CreateTable(&structs.User{})
}

func CreateTable(models interface{}) {
	db.AutoMigrate(models)
}

func Create(value interface{}) bool {
	d := db.Create(value)
	err := d.Error
	if err != nil {
		return false
	} else {
		fmt.Printf("d.RowsAffected: %v\n", d.RowsAffected)
		return true
	}
}

func CheckUserExist(username string, models interface{}) bool {
	d := db.Where("name=?", username).Limit(1).Find(models)
	if d.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func LoginCheck(username string) (bool, string) {
	result := false
	user := structs.User{}
	d := db.Where("name=?", username).Limit(1).Find(&user)
	if d.RowsAffected > 0 {
		result = true
	} else {
		result = false
	}

	if result {
		token, err := tokentool.GenerateToken(user.ID)
		if err != nil {
			return false, ""
		} else {
			return true, token
		}

	} else {
		return false, ""
	}
}

func GetUserByID(id int, models interface{}) {
	db.Find(models, id)
}

// 根据选择的字段插入
func CreateBySelectFields(fields []string, value interface{}) {
	db.Select(fields).Create(value)
}

// 批量插入
func CreateMulti(value interface{}) {
	d := db.Create(value)
	fmt.Printf("d.RowsAffected: %v\n", d.RowsAffected)
}

// 批量插入
func CreateBatch(value interface{}) {
	d := db.CreateInBatches(value, 100)
	fmt.Printf("d.RowsAffected: %v\n", d.RowsAffected)
}

// 忽略钩子
func CreateIgnoreHook(value interface{}) {
	db.Session(&gorm.Session{SkipHooks: true}).Create(value)
}

// 通过map创建
func CreateByMap(value map[string]interface{}, model interface{}) {
	db.Model(model).Create(value)
}

// 冲突啥也不做
func DoNothingWhenInsert(value interface{}) {
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(value)
}

// 通过主键查找
func FindByKey(value interface{}) {
	db.Find(value, 1)
}

// 查找后映射到map类型
func FindAndConvertToMap(value interface{}, result *map[string]interface{}) {
	db.Model(value).First(result)
}

// 多主键查找
func FindByMultiKey(value interface{}, keys []int) {
	db.Find(value, keys)
}

// 条件查找
func FindByCondition(value interface{}, condition map[string]interface{}) {
	// db.Where("name=?", "haha").Order("id desc").Find((value))
	db.Where(condition).Order("id desc").Find(value)
}

// struct 查找，会忽略0值，例如传入ID=0，这个ID=0不会作为查询条件
func FindByStruct(value interface{}, condition interface{}) {
	db.Where(condition).Find(value)
}

// 指定结构体查询字段
func FindByStructWithReturnFields(value interface{}, condition interface{}, returnFields ...string) {
	db.Where(condition, condition).Order("id desc").Find(value)
}
