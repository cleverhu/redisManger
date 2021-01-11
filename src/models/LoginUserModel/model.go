package LoginUserModel

type LoginUser struct {
	ID         int64  `json:"uid" gorm:"column:id;primary_key"`
	Username   string `json:"username" gorm:"column:username"`
	Password   string `json:"password" gorm:"column:password"`
	UpdateTime string `json:"-" gorm:"column:update_time;"`
}
