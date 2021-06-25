package model

type Subscription struct {
	ID                 int      `gorm:"primarykey"`
	CreatedBy          string   `gorm:"column:created_by"`
	CreatedTS          string   `gorm:"column:created_ts"`
	UpdatedBy          string   `gorm:"column:updated_by"`
	UpdatedTS          string   `gorm:"column:updated_ts"`
	CustomerID         int      `gorm:"column:customer_id"`
	Customer           Customer `gorm:"foreignKey:customer_id"`
	ExternalID         string   `gorm:"column:external_id"`
	InstanceType       string   `gorm:"column:instance_type"`
	SubscriptionStatus string   `gorm:"column:subscription_status"`
}

func (Subscription) TableName() string {
	return "subscription"
}
