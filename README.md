# go-ginapi

1. src下新建文件夹

```go
mkdir  ginapi
```
2. init mod

```go
go mod init ginapi
```
3. 安装所有依赖包

```go
go get -u gorm.io/gen //orm 框架
// gin 框架
go get -u github.com/gin-gonic/gin
//jwt 框架
go get -u github.com/dgrijalva/jwt-go
//yaml解析
go get gopkg.in/yaml.v3
//判断slice是否包含某个数值
go get golang.org/x/exp/slices
// s3 sdk
go get github.com/aws/aws-sdk-go/service/s3
// elasticsearch sdk
go get -u github.com/elastic/go-elasticsearch/v7
// elasticsearch dsl
go get "github.com/aquasecurity/esquery"
//操作json
go get "github.com/tidwall/gjson"

go get go.mongodb.org/mongo-driver/mongo
```
4. 项目结构介绍
- controllers 包包含了所有路由对应的函数
- utils 包是jwt生成以及验证的包
- mysqlexample 是MySQL orm CRUD
- middlewares 是JWT 验证中间件
- structs 定义了各种结构体

待补充：
中间价刷新token
权限控制
自定义err
丰富支持的数据库
支持SSO
