package common

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"redisManger/src/dbs"
	"strings"
	"time"
)

var orm *gorm.DB

const (
	USER = iota
	REDIS
)

func init() {
	orm = dbs.Orm
}

func Cors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true                                                                                                 //允许所有域名
	config.AllowMethods = []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"}                                                     //允许请求的方法
	config.AllowHeaders = []string{"tus-resumable", "upload-length", "upload-metadata", "cache-control", "x-requested-with", "*"} //允许的Header
	return cors.New(config)
}

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		var uid int64

		path := context.Request.RequestURI
		method := context.Request.Method
		if path != "/user/login" && path != "/user/check" && strings.Index(path, "/files") == -1 {
			token := context.Request.Header.Get("Token")
			_, claims, err := ParseToken(token)
			if err != nil {
				context.JSON(http.StatusUnauthorized, gin.H{"message": "token error"})
				context.Abort()
			}

			uid = claims.UserId

		}
		if strings.Index(path, "logs") != -1 {
			goto LABEL
		}
		if strings.Index(path, "user") != -1 {
			//user := LoginUserModel.LoginUser{}
			//_ = context.ShouldBindJSON(&user)
			//orm.Table("users").Select("id").Where("username", user.Username).First(&uid)
			//if uid != 0 {
			//	orm.Table("user_log").Exec("insert into user_log (uid,aid,path,method,logtime) values (?,?,?,?,?)", uid, USER, path, method, time.Now())
			//}
		} else {
			orm.Table("user_log").Exec("insert into user_log (uid,aid,path,method,logtime) values (?,?,?,?,?)", uid, REDIS, path, method, time.Now())
		}
	LABEL:
		context.Next()
	}
}

func Logs() gin.HandlerFunc {
	return func(context *gin.Context) {

		context.Next()
	}
}
