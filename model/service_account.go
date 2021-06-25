package model

type ServiceAccount struct {
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

func (ServiceAccount) TableName() string {
	return "service_accounts"
}
