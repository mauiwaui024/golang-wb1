package handler

import (
	"net/http"
	"wb-1/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.GET("/", h.getHomePage)

	router.GET("/order/:id", h.getOrderById)
	///написать функцию считывания из файла

	//ахуенно

	router.NoRoute(func(c *gin.Context) {
		// Render the 404 page HTML template
		c.HTML(http.StatusNotFound, "404.html", nil)
	})

	return router
}
