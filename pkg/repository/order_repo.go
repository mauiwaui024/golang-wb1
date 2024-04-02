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

func (r *Repository) GetAllOrdersFromDB() ([]golangwb1.Order, error) {
	query := `
	SELECT
    o.order_uid,
    o.track_number,
    o.entry,
    o.locale,
    o.internal_signature,
    o.customer_id,
    o.delivery_service,
    o.shardkey,
    o.sm_id,
    o.date_created,
    o.oof_shard,
    d.name AS delivery_name,
    d.phone AS delivery_phone,
    d.zip AS delivery_zip,
    d.city AS delivery_city,
    d.address AS delivery_address,
    d.region AS delivery_region,
    d.email AS delivery_email,
    p.transaction AS payment_transaction,
    p.request_id AS payment_request_id,
    p.currency AS payment_currency,
    p.provider AS payment_provider,
    p.amount AS payment_amount,
    p.payment_dt AS payment_payment_dt,
    p.bank AS payment_bank,
    p.delivery_cost AS payment_delivery_cost,
    p.goods_total AS payment_goods_total,
    p.custom_fee AS payment_custom_fee,
    oi.chrt_id AS item_chrt_id,
    oi.track_number AS item_track_number,
    oi.price AS item_price,
    oi.rid AS item_rid,
    oi.name AS item_name,
    oi.sale AS item_sale,
    oi.size AS item_size,
    oi.total_price AS item_total_price,
    oi.nm_id AS item_nm_id,
    oi.brand AS item_brand,
    oi.status AS item_status
FROM
    orders o
LEFT JOIN
    delivery d ON o.order_uid = d.order_uid
LEFT JOIN
    payment p ON o.order_uid = p.order_uid
LEFT JOIN
    order_items oi ON o.order_uid = oi.order_uid
	`

	// Выполнить запрос SQL
	rows, err := r.DataBase.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Создать слайс для хранения заказов
	var orders []golangwb1.Order
	var orderItems []golangwb1.OrderItem

	// Итерироваться по результатам запроса и создавать структуры Order
	for rows.Next() {
		var order golangwb1.Order
		var delivery golangwb1.DeliveryInfo
		var payment golangwb1.PaymentInfo
		var orderItem golangwb1.OrderItem

		// var orderItems []golangwb1.OrderItem

		// Сканировать данные из результата запроса в переменные
		err := rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SMID,
			&order.DateCreated,
			&order.OOFShard,
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
			&payment.Transaction,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDT,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
			&orderItem.ChrtID,
			&orderItem.TrackNumber,
			&orderItem.Price,
			&orderItem.RID,
			&orderItem.Name,
			&orderItem.Sale,
			&orderItem.Size,
			&orderItem.TotalPrice,
			&orderItem.NMID,
			&orderItem.Brand,
			&orderItem.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		order.Delivery = delivery
		order.Payment = payment

		// order.Items = append(order.Items, orderItem)
		orderItems = append(orderItems, orderItem)

		// Назначить слайс элементов заказа текущему заказу
		order.Items = orderItems

		fmt.Println(order.Items)
		// Добавить заказ в слайс
		orders = append(orders, order)
	}

	// Проверить наличие ошибок после итерации по результатам запроса
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %v", err)
	}

	return orders, nil
}
