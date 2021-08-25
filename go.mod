module redisManger

go 1.14

require (
	github.com/360EntSecGroup-Skylar/excelize/v2 v2.3.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/jinzhu/gorm v1.9.16
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.4 // indirect
	golang.org/x/text v0.3.3
	gopkg.in/yaml.v2 v2.3.0
)

replace gorm.io/gorm v1.23.3 => github.com/go-gorm/gorm v1.23.3
