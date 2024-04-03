package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
	golangwb1 "wb-1"

	"github.com/icrowley/fake"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := InitConfig(); err != nil {
		logrus.Fatal("error initializing configs", err.Error())
	}

	clusterID := viper.GetString("nats.cluster_id")
	natsURL := viper.GetString("nats.url")

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	sc, err := stan.Connect(clusterID, "client_2", stan.NatsConn(nc))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	// генерим заказы и отправляем в канал, принтуем UID каждого заказа
	for i := 0; i < 3000; i++ {
		order := createFakeOrder()
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Fatal(err)
		}
		channel := "orders"
		err = sc.Publish(channel, orderJSON)
		if err != nil {
			log.Fatal(err)
		}
	}

	//Отправляем в канал не валидные данные и получаем ошибку
	// type OrderNotValid struct {
	// 	OrderID   string
	// 	ProductID string
	// 	Quantity  int
	// }

	// var order OrderNotValid = OrderNotValid{OrderID: "2", ProductID: "432", Quantity: 11}
	// NotValidOrderJson, err := json.Marshal(order)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// channel := "orders"
	// err = sc.Publish(channel, NotValidOrderJson)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	log.Printf("Order data sent to channel")

}

func createFakeOrder() golangwb1.Order {

	order := golangwb1.Order{
		OrderUID:          fake.DigitsN(13),
		TrackNumber:       fake.CharactersN(10),
		Entry:             fake.CharactersN(20),
		Locale:            "en",
		InternalSignature: fake.CharactersN(15),
		CustomerID:        fake.CharactersN(8),
		DeliveryService:   fake.Company(),
		ShardKey:          fake.CharactersN(10),
		SMID:              fake.Day(),
		DateCreated:       "2021-11-26T06:22:19Z",
		OOFShard:          fake.CharactersN(5),
	}

	order.Delivery = golangwb1.DeliveryInfo{
		Name:    fake.FullName(),
		Phone:   fake.Phone(),
		Zip:     fake.Zip(),
		City:    fake.City(),
		Address: fake.StreetAddress(),
		Region:  fake.State(),
		Email:   fake.EmailAddress(),
	}

	order.Payment = golangwb1.PaymentInfo{
		Transaction:  fake.CharactersN(10),
		RequestID:    fake.CharactersN(8),
		Currency:     fake.CurrencyCode(),
		Provider:     fake.Company(),
		Amount:       3,
		PaymentDT:    time.Now().Unix(),
		Bank:         "Tinkoff",
		DeliveryCost: rand.Intn(10000),

		GoodsTotal: rand.Intn(100),
		CustomFee:  rand.Intn(10),
	}
	for i := 0; i < 3; i++ {
		item := golangwb1.OrderItem{
			ChrtID:      i + 1,
			TrackNumber: fake.CharactersN(8),
			Price:       rand.Intn(30000),
			RID:         fake.CharactersN(5),
			Name:        fake.ProductName(),
			Sale:        rand.Intn(30),
			Size:        "0",
			TotalPrice:  rand.Intn(1000),
			NMID:        rand.Intn(599),
			Brand:       fake.Brand(),
			Status:      rand.Intn(400),
		}
		order.Items = append(order.Items, item)
	}
	// fmt.Println("orderUID of generated order ", order.OrderUID)
	return order
}
