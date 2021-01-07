package dbs

import (
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

var (
	Rds *redis.Client
	Orm *gorm.DB
)

func init() {
	err := fmt.Errorf("")
	Orm, err = gorm.Open("mysql", "root:123456@tcp(101.132.107.3:3306)/rbac?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		log.Fatal(err)
	}

	mysqlDB := Orm.DB()
	mysqlDB.SetConnMaxLifetime(30 * time.Second)
	mysqlDB.SetMaxIdleConns(10)
	mysqlDB.SetMaxIdleConns(5)
	Orm.LogMode(true)

	//初始化redis
	Rds = redis.NewClient(&redis.Options{
		Addr:     "101.132.107.3:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
