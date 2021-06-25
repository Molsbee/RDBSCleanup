package model

type Customer struct {
	ID        uint   `gorm:"primarykey"`
	CreatedBy string `gorm:"column:created_by"`
	CreatedTS string `gorm:"column:created_ts"`
	UpdatedBy string `gorm:"column:updated_by"`
	Alias     string `gorm:"column:alias"`
}

func (Customer) TableName() string {
	return "customer"
}
