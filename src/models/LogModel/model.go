package LogModel

import "redisManger/src/dbs"

type LogImpl struct {
	ID      int64  `json:"id" gorm:"column:id;primary_key"`
	Uid     int64  `json:"uid" gorm:"column:uid"`
	Path    string `json:"path" gorm:"column:path"`
	Method  string `json:"method" gorm:"column:method"`
	LogTime string `json:"log_time" gorm:"column:logtime"`
	// only for query.
	Action string `json:"action" gorm:"column:action"`
	Aid    int64  `json:"aid" gorm:"column:aid"`
}

func (*LogImpl) TableName() string {
	return "user_log"
}

func init() {
	dbs.Orm.AutoMigrate(&LogImpl{})
}
