package dbs

import (
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type redisConf struct {
	Host string `yaml:"host"`
	Post int    `yaml:"port"`
}

type T struct {
	Redis redisConf `yaml:"redis"`
}

var (
	Rds *redis.Client
	Orm *gorm.DB
)

func init() {
	err := fmt.Errorf("")
	Orm, err = gorm.Open("mysql", "root:123456@tcp(101.132.107.3:3306)/redis_manger?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	mysqlDB := Orm.DB()
	mysqlDB.SetConnMaxLifetime(30 * time.Second)
	mysqlDB.SetMaxIdleConns(10)
	mysqlDB.SetMaxIdleConns(5)
	Orm.LogMode(true)

	//初始化redis
	t := &T{}
	file, err := ioutil.ReadFile("./config.yaml")
	workDir, _ := os.Getwd()
	fmt.Println("work_dir:", workDir)
	fmt.Println("file_err", err)
	fmt.Println("file_content", string(file))
	err = yaml.Unmarshal(file, t)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println(t)
	//redis.Options{
	//	Network:            "",
	//	Addr:               "",
	//	Dialer:             nil,
	//	OnConnect:          nil,
	//	Password:           "",
	//	DB:                 0,
	//	MaxRetries:         0,
	//	MinRetryBackoff:    0,
	//	MaxRetryBackoff:    0,
	//	DialTimeout:        0,
	//	ReadTimeout:        0,
	//	WriteTimeout:       0,
	//	PoolSize:           0,
	//	MinIdleConns:       0,
	//	MaxConnAge:         0,
	//	PoolTimeout:        0,
	//	IdleTimeout:        0,
	//	IdleCheckFrequency: 0,
	//	TLSConfig:          nil,
	//}
	Rds = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", "101.132.107.3", 6379),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println(Rds)
}
