package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) getHomePage(c *gin.Context) {

}

func (h *Handler) getOrderById(c *gin.Context) {
	//		orderID := c.Param("id")
	//		fmt.Println("pam pam pam", orderID)
	//		orderID = trimFirstRune(orderID)
	//		orderIDInt, err := strconv.Atoi(orderID)
	//		if err != nil {
	//			c.JSON(http.StatusBadRequest, gin.H{"error": "This order does not exist, sorry"})
	//			return
	//		}
	//		order, err := h.services.GetOrderById(orderIDInt)
	//		if err != nil {
	//			c.JSON(http.StatusBadRequest, gin.H{"error": "This order does not exist in this database"})
	//		}
	//		c.HTML(http.StatusOK, "single_order.html", gin.H{
	//			"order": order,
	//		})
	//	}
	//
	//	func trimFirstRune(s string) string {
	//		for i := range s {
	//			if i > 0 {
	//				return s[i:]
	//			}
	//		}
	//		return ""
}
