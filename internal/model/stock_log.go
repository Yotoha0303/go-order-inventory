package model

import "time"

const (
	StockBizInit      = 1
	StockBizManualAdd = 2
)

type StockLog struct {
	ID             int64     `gorm:"primaryKey;autoIncrement;type:bigint" json:"id"`
	ProductID      int64     `gorm:"type:bigint;not null;index" json:"product_id"`
	ChangeQuantity int64     `gorm:"type:bigint;not null" json:"change_quantity"`
	BeforeQuantity int64     `gorm:"type:bigint;not null" json:"before_quantity"`
	AfterQuantity  int64     `gorm:"type:bigint;not null" json:"after_quantity"`
	BizType        int8      `gorm:"type:tinyint;not null;index" json:"biz_type"`
	BizID          *int64    `gorm:"type:bigint;index" json:"biz_id"`
	Remark         string    `gorm:"type:varchar(255);not null;default:''" json:"remark"`
	CreatedAt      time.Time `json:"created_at"`
}

func (StockLog) TableName() string {
	return "stock_logs"
}
