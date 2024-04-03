package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	golangwb1 "wb-1"
	"wb-1/pkg/handler"
	"wb-1/pkg/repository"
	"wb-1/pkg/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if err := InitConfig(); err != nil {
		logrus.Fatal("error initializing configs", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatal("error loading env variable", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.db"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatal("Failed to initialize db", err.Error())
	}

	//сначала создаем репозиторий, потом сервис, который зависит от репозитория,
	//а потом хэндлер, который зависит от сервисов
	repos := repository.NewRepository(db)
	services := service.NewService(repos)

	//подключаемся к нэтс-стриминг
	sc, err := connectToNATSStreaming()
	if err != nil {
		logrus.Fatal("Failed to connect to NATS Streaming server: ", err.Error())
	}
	defer sc.Close()

	//забираем всё заказы из бд
	ordersForCache, err := services.GetAllOrdersFromDB()
	if err != nil {
		logrus.Fatal("Failed to get all orders from db", err.Error())
	}
	//создаем кэш
	cache := golangwb1.NewCache()
	// Восстанавливаем данные из базы данных
	cache.RestoreFromDatabase(ordersForCache)
	fmt.Println("Contents of the cache map after it is loaded from db:")
	for _, value := range cache.Orders {
		fmt.Println(value)
	}
	// Создание экземпляра хендлера
	handlers := handler.NewHandler(services, sc, cache)

	//подписываемся на канал
	go func() {
		err = handlers.SubscribeToChannel("orders")
		if err != nil {
			logrus.Fatal("Failed to subscribe to channel: ", err)
		}
	}()

	srv := new(golangwb1.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatal(err)
		}
	}()
	logrus.Print("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("APP shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error occurred on db shutting down: %s", err.Error())
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func connectToNATSStreaming() (stan.Conn, error) {
	clusterID := viper.GetString("nats.cluster_id")
	clientID := viper.GetString("nats.client_id")
	natsURL := viper.GetString("nats.url")

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		return nil, err
	}

	return sc, nil
}
