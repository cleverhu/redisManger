package LogModel

type LogImpl struct {
	Uid     int64  `json:"uid" gorm:"column:uid"`
	Action  string `json:"action" gorm:"column:action"`
	Path    string `json:"path" gorm:"column:path"`
	Method  string `json:"method" gorm:"column:method"`
	LogTime string `json:"logtime" gorm:"column:logtime"`
}


