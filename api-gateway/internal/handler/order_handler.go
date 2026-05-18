package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	orderpb "sneaker-store/order-service/proto"
)

type OrderHandler struct {
	client orderpb.OrderServiceClient
}

func NewOrderHandler(client orderpb.OrderServiceClient) *OrderHandler {
	return &OrderHandler{client: client}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req orderpb.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.client.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.client.GetOrder(c.Request.Context(), &orderpb.GetOrderRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	resp, err := h.client.ListOrders(c.Request.Context(), &orderpb.ListOrdersRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.client.UpdateOrderStatus(c.Request.Context(), &orderpb.UpdateOrderStatusRequest{
		Id:     id,
		Status: body.Status,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.client.CancelOrder(c.Request.Context(), &orderpb.CancelOrderRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) GetOrdersByUser(c *gin.Context) {
	userID := c.Param("user_id")
	resp, err := h.client.GetOrdersByUser(c.Request.Context(), &orderpb.GetOrdersByUserRequest{UserId: userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) GetOrdersByStatus(c *gin.Context) {
	orderStatus := c.Param("status")
	resp, err := h.client.GetOrdersByStatus(c.Request.Context(), &orderpb.GetOrdersByStatusRequest{Status: orderStatus})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) GetOrderStats(c *gin.Context) {
	resp, err := h.client.GetOrderStats(c.Request.Context(), &orderpb.GetOrderStatsRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) GetOrderItems(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.client.GetOrderItems(c.Request.Context(), &orderpb.GetOrderItemsRequest{OrderId: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) GetTotalRevenue(c *gin.Context) {
	resp, err := h.client.GetTotalRevenue(c.Request.Context(), &orderpb.GetTotalRevenueRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) GetOrdersByDateRange(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	resp, err := h.client.GetOrdersByDateRange(c.Request.Context(), &orderpb.GetOrdersByDateRangeRequest{From: from, To: to})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) CountOrdersByUser(c *gin.Context) {
	userID := c.Param("user_id")
	resp, err := h.client.CountOrdersByUser(c.Request.Context(), &orderpb.CountOrdersByUserRequest{UserId: userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
