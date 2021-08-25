package dbs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"
)

type redisConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

type mysqlConf struct {
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DbName   string `yaml:"dbname"`
}

type T struct {
	Redis redisConf `yaml:"redis"`
	Mysql mysqlConf `yaml:"mysql"`
}

var (
	Rds *redis.Client
	Orm *gorm.DB
)

func init() {
	t := &T{}
	file, err := ioutil.ReadFile("./config.yaml")
	err = yaml.Unmarshal(file, t)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	Orm, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local", t.Mysql.UserName, t.Mysql.Password, t.Mysql.Host, t.Mysql.Port, t.Mysql.DbName))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	mysqlDB := Orm.DB()
	mysqlDB.SetConnMaxLifetime(30 * time.Second)
	mysqlDB.SetMaxIdleConns(10)
	mysqlDB.SetMaxIdleConns(5)
	Orm.LogMode(true)
	Rds = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", t.Redis.Host, t.Redis.Port),
		Password: t.Redis.Password, // no password set
		DB:       t.Redis.Db,       // use default DB
	})
	_, err = Rds.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
}
