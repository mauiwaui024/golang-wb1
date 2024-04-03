package handler

import (
	"encoding/json"
	"net/http"
	golangwb1 "wb-1"
	"wb-1/pkg/service"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	services *service.Service
	stanConn stan.Conn
	cache    *golangwb1.Cache
}

func NewHandler(services *service.Service, stanConn stan.Conn, cache *golangwb1.Cache) *Handler {
	return &Handler{
		services: services,
		stanConn: stanConn,
		cache:    cache,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("./cmd/templates/*")
	router.GET("/", h.getForm)

	router.POST("/order", h.getOrderById)

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
	})
	return router
}

func (h *Handler) SubscribeToChannel(channel string) error {
	_, err := h.stanConn.Subscribe(channel, func(msg *stan.Msg) {
		// Обработка полученного сообщения
		// fmt.Printf("%s\n", msg.Data)
		var completeOrder golangwb1.Order
		err := json.Unmarshal(msg.Data, &completeOrder)
		if err != nil {
			if e, ok := err.(*json.UnmarshalTypeError); ok {
				logrus.WithError(err).Errorf("Failed to unmarshal complete order data: expected type %s, got %s", e.Type, e.Value)
			} else {
				logrus.WithError(err).Error("Failed to unmarshal complete order data")
			}
			return
		}
		// fmt.Println(completeOrder.OrderUID)
		//чекаем заказ на валидность
		if !isValidOrder(completeOrder) {
			logrus.Error("Invalid order structure")
			return
		}
		//запихиваем в бд
		err = h.services.CreateOrder(completeOrder)
		if err != nil {
			logrus.WithError(err).Error("Failed to create order in db")
			return
		}
		//если в бд заказ записался, добавляем его и в кэш
		h.cache.AddOrder(completeOrder)
		logrus.Printf("Order with orderUID = %s added to db and cache", completeOrder.OrderUID)

	}, stan.DurableName(channel))
	if err != nil {
		return err
	}
	return nil
}

func isValidOrder(order golangwb1.Order) bool {
	if order.OrderUID == "" || order.TrackNumber == "" || order.CustomerID == "" {
		return false
	}
	if order.Delivery.Name == "" || order.Delivery.Phone == "" {
		return false
	}
	if order.Payment.Transaction == "" {
		return false
	}
	for _, item := range order.Items {
		if item.ChrtID == 0 || item.TrackNumber == "" || item.Price == 0 || item.Name == "" {
			return false
		}
	}
	return true
}

// service.HandleMessage(msg.Data)

// fmt.Println(string(msg.Data))
// fmt.Println(completeOrder.Payment)
// fmt.Println(completeOrder)
// fmt.Println(completeOrder.OrderItems[0])
