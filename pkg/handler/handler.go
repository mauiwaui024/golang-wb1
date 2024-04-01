package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	golangwb1 "wb-1"
	"wb-1/pkg/service"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/stan.go"
)

type Handler struct {
	services *service.Service
	stanConn stan.Conn // внедряем зависимость
}

func NewHandler(services *service.Service, stanConn stan.Conn) *Handler {
	return &Handler{
		services: services,
		stanConn: stanConn, // Присваиваем переданное подключение к полю структуры
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.GET("/", h.getHomePage)

	router.GET("/order/:id", h.getOrderById)

	router.NoRoute(func(c *gin.Context) {
		// Render the 404 page HTML template
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
		/////
		// service.HandleMessage(msg.Data)

		fmt.Println(string(msg.Data))
		// fmt.Println(completeOrder.Payment)
		// fmt.Println(completeOrder)
		// fmt.Println(completeOrder.OrderItems[0])
		////////
		if err != nil {
			fmt.Println("Failed to unmarshal complete order data:", err)
			return
		}
		err = h.services.CreateOrder(completeOrder)
		if err != nil {
			fmt.Println("Failed to create new instance of ORDER in DB", err)
		}

	}, stan.DurableName(channel))
	if err != nil {
		return err
	}
	return nil
}
