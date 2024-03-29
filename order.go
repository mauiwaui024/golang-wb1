package golangwb1

import "time"

type Order struct {
	OrderUID    string    `json:"order_uid"`
	TrackNumber string    `json:"track_number"`
	DateCreated time.Time `json:"date_created"`
	ShardKey    string    `json:"shardkey"`
	SMID        int       `json:"sm_id"`
	OOFShard    string    `json:"oof_shard"`
}

// DeliveryInfo represents the delivery information
type DeliveryInfo struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

// PaymentInfo represents the payment information
type PaymentInfo struct {
	Transaction  string    `json:"transaction"`
	RequestID    string    `json:"request_id"`
	Currency     string    `json:"currency"`
	Provider     string    `json:"provider"`
	Amount       int       `json:"amount"`
	PaymentDT    time.Time `json:"payment_dt"`
	Bank         string    `json:"bank"`
	DeliveryCost int       `json:"delivery_cost"`
	GoodsTotal   int       `json:"goods_total"`
	CustomFee    int       `json:"custom_fee"`
}

// CompleteOrder represents the complete order data including Order, DeliveryInfo, and PaymentInfo
type CompleteOrder struct {
	Order    Order        `json:"order"`
	Delivery DeliveryInfo `json:"delivery_info"`
	Payment  PaymentInfo  `json:"payment_info"`
}
