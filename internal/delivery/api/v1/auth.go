package v1

import (
	"net/http"

	"github.com/zsandibe/medods-service/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) Login(c *gin.Context) {
	var loginRequest domain.LoginRequest
	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if loginRequest.Guid == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid GUID format"})
		return
	}

	tokenPair, err := h.service.Create(loginRequest.Guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tokenPair)
}

func (h *Handler) Refresh(c *gin.Context) {
	var refreshRequest domain.RefreshRequest
	if err := c.BindJSON(&refreshRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if refreshRequest.SessionID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SessionID format"})
		return
	}

	refreshedTokenPair, err := h.service.Update(refreshRequest.SessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, refreshedTokenPair)
}

func (h *Handler) GetAllSessions(c *gin.Context) {
	sessions, err := h.service.GetAllSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(sessions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sessions not found"})
		return
	}
	c.JSON(http.StatusOK, sessions)
}
