package models

type StockDecreaseLog struct {
	ID           int64 `json:"id" gorm:"primaryKey"`
	OrderId      int64 `json:"order_id"`
	ProductRefer int64 `json:"product_id"`
}
