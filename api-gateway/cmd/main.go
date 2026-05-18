package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	gatewayclient "sneaker-store/api-gateway/internal/client"
	"sneaker-store/api-gateway/internal/handler"
	"sneaker-store/api-gateway/internal/middleware"
	orderpb "sneaker-store/order-service/proto"
	productpb "sneaker-store/product-service/proto"
	userpb "sneaker-store/user-service/proto"
)

func main() {
	productAddr := envOr("PRODUCT_SERVICE_ADDR", "localhost:50051")
	orderAddr := envOr("ORDER_SERVICE_ADDR", "localhost:50052")
	userAddr := envOr("USER_SERVICE_ADDR", "localhost:50053")
	aiAddr := envOr("AI_SERVICE_ADDR", "localhost:50054")
	httpPort := envOr("HTTP_PORT", "8080")

	clients, err := gatewayclient.NewClients(productAddr, orderAddr, userAddr, aiAddr)
	if err != nil {
		log.Fatalf("init clients: %v", err)
	}

	productH := handler.NewProductHandler(clients.Product)
	orderH := handler.NewOrderHandler(clients.Order)
	userH := handler.NewUserHandler(clients.User)
	aiH := handler.NewAIHandler(clients.AI)

	r := gin.Default()

	// CORS for Next.js frontend
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// ── Public routes ──────────────────────────────────────────────────────────
	products := r.Group("/api/v1/products")
	{
		products.POST("", productH.CreateProduct)
		products.GET("", productH.ListProducts)
		products.GET("/search", productH.SearchProducts)
		products.GET("/brands", productH.GetBrands)
		products.GET("/stats", productH.GetProductStats)
		products.GET("/low-stock", productH.GetLowStockProducts)
		products.GET("/by-brand/:brand", productH.GetProductsByBrand)
		products.POST("/bulk-delete", productH.BulkDeleteProducts)
		products.GET("/:id", productH.GetProduct)
		products.PUT("/:id", productH.UpdateProduct)
		products.DELETE("/:id", productH.DeleteProduct)
		products.PATCH("/:id/stock", productH.UpdateStock)
	}

	orders := r.Group("/api/v1/orders")
	{
		orders.POST("", orderH.CreateOrder)
		orders.GET("", orderH.ListOrders)
		orders.GET("/stats", orderH.GetOrderStats)
		orders.GET("/revenue", orderH.GetTotalRevenue)
		orders.GET("/by-status/:status", orderH.GetOrdersByStatus)
		orders.GET("/by-date", orderH.GetOrdersByDateRange)
		orders.GET("/user/:user_id", orderH.GetOrdersByUser)
		orders.GET("/user/:user_id/count", orderH.CountOrdersByUser)
		orders.GET("/:id", orderH.GetOrder)
		orders.GET("/:id/items", orderH.GetOrderItems)
		orders.PATCH("/:id/status", orderH.UpdateOrderStatus)
		orders.DELETE("/:id/cancel", orderH.CancelOrder)
	}

	users := r.Group("/api/v1/users")
	{
		users.POST("", userH.CreateUser)
		users.POST("/login", userH.Login)
		users.GET("", userH.ListUsers)
		users.GET("/search", userH.SearchUsers)
		users.GET("/by-email", userH.GetUserByEmail)
		users.GET("/stats", userH.GetUserStats)
		users.GET("/:id", userH.GetUser)
		users.PUT("/:id", userH.UpdateUser)
		users.DELETE("/:id", userH.DeleteUser)
		users.PATCH("/:id/password", userH.ChangePassword)
		users.PATCH("/:id/password/reset", userH.ResetPassword)
		users.PATCH("/:id/status", userH.UpdateUserStatus)
	}

	// ── AI routes ─────────────────────────────────────────────────────────────
	aiRoutes := r.Group("/api/v1/ai")
	{
		aiRoutes.POST("/chat", aiH.Chat)
		aiRoutes.POST("/recommend", aiH.Recommend)
		aiRoutes.POST("/search-by-style", aiH.SearchByStyle)
		aiRoutes.GET("/trending", aiH.Trending)
	}

	// ── Admin routes (JWT + ADMIN role required) ───────────────────────────────
	admin := r.Group("/api/v1/admin", middleware.RequireAuth(), middleware.RequireAdmin())
	{
		// Products
		admin.POST("/products", productH.CreateProduct)
		admin.GET("/products", productH.ListProducts)
		admin.GET("/products/low-stock", productH.GetLowStockProducts)
		admin.GET("/products/stats", productH.GetProductStats)
		admin.GET("/products/:id", productH.GetProduct)
		admin.PUT("/products/:id", productH.UpdateProduct)
		admin.DELETE("/products/:id", productH.DeleteProduct)
		admin.PATCH("/products/:id/stock", productH.UpdateStock)
		admin.POST("/products/bulk-delete", productH.BulkDeleteProducts)

		// Orders
		admin.GET("/orders", orderH.ListOrders)
		admin.GET("/orders/stats", orderH.GetOrderStats)
		admin.GET("/orders/revenue", orderH.GetTotalRevenue)
		admin.GET("/orders/by-status/:status", orderH.GetOrdersByStatus)
		admin.GET("/orders/:id", orderH.GetOrder)
		admin.GET("/orders/:id/items", orderH.GetOrderItems)
		admin.PATCH("/orders/:id/status", orderH.UpdateOrderStatus)
		admin.DELETE("/orders/:id/cancel", orderH.CancelOrder)

		// Users
		admin.GET("/users", userH.ListUsers)
		admin.GET("/users/search", userH.SearchUsers)
		admin.GET("/users/stats", userH.GetUserStats)
		admin.GET("/users/:id", userH.GetUser)
		admin.DELETE("/users/:id", userH.DeleteUser)
		admin.PATCH("/users/:id/status", userH.UpdateUserStatus)
		admin.PATCH("/users/:id/password/reset", userH.ResetPassword)

		// Aggregated stats dashboard endpoint
		admin.GET("/stats", func(c *gin.Context) {
			ctx := context.Background()
			pStats, _ := clients.Product.GetProductStats(ctx, &productpb.GetProductStatsRequest{})
			oStats, _ := clients.Order.GetOrderStats(ctx, &orderpb.GetOrderStatsRequest{})
			uStats, _ := clients.User.GetUserStats(ctx, &userpb.GetUserStatsRequest{})
			rev, _ := clients.Order.GetTotalRevenue(ctx, &orderpb.GetTotalRevenueRequest{})

			type DashStats struct {
				TotalProducts int32   `json:"total_products"`
				TotalOrders   int32   `json:"total_orders"`
				TotalUsers    int32   `json:"total_users"`
				TotalRevenue  float64 `json:"total_revenue"`
				PendingOrders int32   `json:"pending_orders"`
				ActiveUsers   int32   `json:"active_users"`
			}
			s := DashStats{}
			if pStats != nil {
				s.TotalProducts = pStats.TotalProducts
			}
			if oStats != nil {
				s.TotalOrders = oStats.TotalOrders
				s.PendingOrders = oStats.PendingOrders
			}
			if uStats != nil {
				s.TotalUsers = uStats.TotalUsers
				s.ActiveUsers = uStats.ActiveUsers
			}
			if rev != nil {
				s.TotalRevenue = rev.Total
			}
			c.JSON(http.StatusOK, s)
		})
	}

	log.Printf("api-gateway listening on :%s", httpPort)
	if err := r.Run(":" + httpPort); err != nil {
		log.Fatalf("run: %v", err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
