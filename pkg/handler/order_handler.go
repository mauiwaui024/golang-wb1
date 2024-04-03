package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getForm(c *gin.Context) {
	c.HTML(http.StatusOK, "form.html", nil)
}

func (h *Handler) getOrderById(c *gin.Context) {
	orderID := c.PostForm("id")
	// fmt.Println("Orderid is", orderID)
	order, found := h.cache.GetOrderById(orderID)
	if !found {
		c.Redirect(http.StatusFound, "/?not_found")
		return
	}
	// for item := range order.Items {
	// 	fmt.Println(order.Items)
	// }
	// fmt.Println("razrzzzzzzzzzzzzz", order.Items)

	c.HTML(http.StatusOK, "order.html", gin.H{
		"OrderUID":          order.OrderUID,
		"TrackNumber":       order.TrackNumber,
		"Entry":             order.Entry,
		"Delivery":          order.Delivery,
		"Payment":           order.Payment,
		"OrderItems":        order.Items,
		"Locale":            order.Locale,
		"InternalSignature": order.InternalSignature,
		"CustomerID":        order.CustomerID,
		"DeliveryService":   order.DeliveryService,
		"ShardKey":          order.ShardKey,
		"SMID":              order.SMID,
		"DateCreated":       order.DateCreated,
		"OOFShard":          order.OOFShard,
	})

	// c.JSON(http.StatusOK, order)
}

// fmt.Println("Contents of the map after adding new order:")
// for _, value := range h.cache.Orders {
// 	fmt.Println(value)
// }
