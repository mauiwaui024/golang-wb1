package main

import (
	"encoding/json"
	"log"
	"time"
	golangwb1 "wb-1"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

func main() {

	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	sc, err := stan.Connect("test-cluster", "client-1", stan.NatsConn(nc))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	order := golangwb1.Order{
		OrderUID:    "123456",
		TrackNumber: "ABC123",
		DateCreated: time.Now(),
		ShardKey:    "shard1",
		SMID:        1,
		OOFShard:    "oof-shard1",
	}

	deliveryInfo := golangwb1.DeliveryInfo{
		Name:    "John Doe",
		Phone:   "123456789",
		Zip:     "12345",
		City:    "New York",
		Address: "123 Main St",
		Region:  "NY",
		Email:   "john@example.com",
	}

	paymentInfo := golangwb1.PaymentInfo{
		Transaction:  "txn123",
		RequestID:    "req123",
		Currency:     "USD",
		Provider:     "PayPal",
		Amount:       100,
		PaymentDT:    time.Now(),
		Bank:         "Bank XYZ",
		DeliveryCost: 10,
		GoodsTotal:   90,
		CustomFee:    5,
	}

	completeOrder := golangwb1.CompleteOrder{
		Order:    order,
		Delivery: deliveryInfo,
		Payment:  paymentInfo,
	}

	// херачим в джсон
	orderData, err := json.Marshal(completeOrder)
	if err != nil {
		log.Fatal(err)
	}

	// и в канал публикуем
	channel := "orders"
	err = sc.Publish(channel, orderData)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Order data sent to channel %s", channel)

	select {}
}
