package LoginUserModel

import "redisManger/src/dbs"

type LoginUser struct {
	ID         int64  `json:"uid" gorm:"column:id;primary_key"`
	Username   string `json:"username" gorm:"column:username"`
	Password   string `json:"password" gorm:"column:password"`
	UpdateTime string `json:"-" gorm:"column:update_time;"`
}

func (u LoginUser) TableName() string {
	return "users"
}

func init() {
	dbs.Orm.AutoMigrate(LoginUser{})
}
