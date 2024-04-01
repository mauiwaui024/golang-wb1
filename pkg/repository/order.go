package repository

import (
	"fmt"
	golangwb1 "wb-1"
)

func (r *Repository) CreateOrder(completeOrder golangwb1.Order) error {
	tx, err := r.DataBase.Begin()
	if err != nil {
		return err
	}

	// 1. Добавить данные в таблицу orders
	fmt.Println("rz")
	createOrdersQuery := `INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err = tx.Exec(createOrdersQuery,
		completeOrder.OrderUID,
		completeOrder.TrackNumber,
		completeOrder.Entry,
		completeOrder.Locale,
		completeOrder.InternalSignature,
		completeOrder.CustomerID,
		completeOrder.DeliveryService,
		completeOrder.ShardKey,
		completeOrder.SMID,
		completeOrder.DateCreated,
		completeOrder.OOFShard)
	if err != nil {
		tx.Rollback()
		return err
	}
	fmt.Println("wwww")
	// 2. Добавить данные в таблицу delivery
	createDeliveryQuery := `INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = tx.Exec(createDeliveryQuery,
		completeOrder.OrderUID,
		completeOrder.Delivery.Name,
		completeOrder.Delivery.Phone,
		completeOrder.Delivery.Zip,
		completeOrder.Delivery.City,
		completeOrder.Delivery.Address,
		completeOrder.Delivery.Region,
		completeOrder.Delivery.Email)
	if err != nil {
		tx.Rollback()
		return err
	}
	// 3. Добавить данные в таблицу payment
	createPaymentQuery := `INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err = tx.Exec(createPaymentQuery,
		completeOrder.OrderUID,
		completeOrder.Payment.Transaction,
		completeOrder.Payment.RequestID,
		completeOrder.Payment.Currency,
		completeOrder.Payment.Provider,
		completeOrder.Payment.Amount,
		completeOrder.Payment.PaymentDT,
		completeOrder.Payment.Bank,
		completeOrder.Payment.DeliveryCost,
		completeOrder.Payment.GoodsTotal,
		completeOrder.Payment.CustomFee)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 4. Добавить данные в таблицу order_items
	createOrderItemsQuery := `INSERT INTO order_items (order_uid, track_number, chrt_id, price, rid, name, sale, size, total_price, nm_id, brand, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	for _, item := range completeOrder.Items {
		_, err := tx.Exec(createOrderItemsQuery,
			completeOrder.OrderUID,
			item.TrackNumber,
			item.ChrtID,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NMID,
			item.Brand,
			item.Status)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
