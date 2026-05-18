package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	productpb "sneaker-store/product-service/proto"
)

type ProductHandler struct {
	client productpb.ProductServiceClient
}

func NewProductHandler(client productpb.ProductServiceClient) *ProductHandler {
	return &ProductHandler{client: client}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req productpb.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.client.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.client.GetProduct(c.Request.Context(), &productpb.GetProductRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	resp, err := h.client.ListProducts(c.Request.Context(), &productpb.ListProductsRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var req productpb.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Id = id
	resp, err := h.client.UpdateProduct(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.client.DeleteProduct(c.Request.Context(), &productpb.DeleteProductRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) SearchProducts(c *gin.Context) {
	req := &productpb.SearchProductsRequest{
		Brand: c.Query("brand"),
		Size:  c.Query("size"),
	}
	resp, err := h.client.SearchProducts(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) UpdateStock(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Delta int32 `json:"delta"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.client.UpdateStock(c.Request.Context(), &productpb.UpdateStockRequest{Id: id, Delta: body.Delta})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) GetProductsByBrand(c *gin.Context) {
	brand := c.Param("brand")
	resp, err := h.client.GetProductsByBrand(c.Request.Context(), &productpb.GetProductsByBrandRequest{Brand: brand})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) GetLowStockProducts(c *gin.Context) {
	threshold := int32(10)
	if t := c.Query("threshold"); t != "" {
		if v, err := strconv.Atoi(t); err == nil {
			threshold = int32(v)
		}
	}
	resp, err := h.client.GetLowStockProducts(c.Request.Context(), &productpb.GetLowStockRequest{Threshold: threshold})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) GetBrands(c *gin.Context) {
	resp, err := h.client.GetBrands(c.Request.Context(), &productpb.GetBrandsRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) GetProductStats(c *gin.Context) {
	resp, err := h.client.GetProductStats(c.Request.Context(), &productpb.GetProductStatsRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) BulkDeleteProducts(c *gin.Context) {
	var body struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.client.BulkDeleteProducts(c.Request.Context(), &productpb.BulkDeleteProductsRequest{Ids: body.IDs})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
