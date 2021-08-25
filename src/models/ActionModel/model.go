package ActionModel

import "redisManger/src/dbs"

type Action struct {
	ID     int64  `json:"id" gorm:"column:id;primary_key"`
	Action string `json:"action" gorm:"column:action"`
}

func (*Action) TableName() string {
	return "actions"
}

func init() {
	dbs.Orm.AutoMigrate(&Action{})
	actions := []*Action{
		{
			ID:     0,
			Action: "USER",
		},
		{
			ID:     1,
			Action: "REDIS",
		},
	}
	for _, action := range actions {
		dbs.Orm.Create(action)
	}
}
