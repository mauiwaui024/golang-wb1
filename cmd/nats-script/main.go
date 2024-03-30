package main

import (
	"fmt"
	"log"
	"os"

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

	filePath := "order.json"

	// Открываем файл
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := os.ReadFile("order.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Parse JSON data into CompleteOrder struct

	// Print the parsed data
	// Similarly, access other fields as needed

	// Call function to complete the order

	// Отправка данных в канал
	channel := "orders"
	err = sc.Publish(channel, data)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Order data sent to channel %s", channel)
}
