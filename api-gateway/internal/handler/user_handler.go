package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sneaker-store/api-gateway/internal/middleware"
	userpb "sneaker-store/user-service/proto"
)

type UserHandler struct {
	client userpb.UserServiceClient
}

func NewUserHandler(client userpb.UserServiceClient) *UserHandler {
	return &UserHandler{client: client}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req userpb.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.client.CreateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.client.GetUser(c.Request.Context(), &userpb.GetUserRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req userpb.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Id = id
	resp, err := h.client.UpdateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.client.DeleteUser(c.Request.Context(), &userpb.DeleteUserRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req userpb.AuthenticateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.client.AuthenticateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !resp.Success {
		c.JSON(http.StatusUnauthorized, gin.H{"error": resp.Message})
		return
	}
	token, err := middleware.GenerateToken(resp.UserId, resp.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user_id": resp.UserId,
		"role":    resp.Role,
		"token":   token,
		"message": resp.Message,
	})
}

func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Query("email")
	resp, err := h.client.GetUserByEmail(c.Request.Context(), &userpb.GetUserByEmailRequest{Email: email})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.client.ChangePassword(c.Request.Context(), &userpb.ChangePasswordRequest{
		Id: id, OldPassword: body.OldPassword, NewPassword: body.NewPassword,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page := int32(1)
	limit := int32(20)
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			page = int32(v)
		}
	}
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = int32(v)
		}
	}
	resp, err := h.client.ListUsers(c.Request.Context(), &userpb.ListUsersRequest{Page: page, Limit: limit})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	resp, err := h.client.SearchUsers(c.Request.Context(), &userpb.SearchUsersRequest{Query: query})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) GetUserStats(c *gin.Context) {
	resp, err := h.client.GetUserStats(c.Request.Context(), &userpb.GetUserStatsRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Active bool `json:"active"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.client.UpdateUserStatus(c.Request.Context(), &userpb.UpdateUserStatusRequest{Id: id, Active: body.Active})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.client.ResetPassword(c.Request.Context(), &userpb.ResetPasswordRequest{Id: id, NewPassword: body.NewPassword})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
