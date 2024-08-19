package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zsandibe/medods-service/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) Login(c *gin.Context) {
	var loginRequest domain.LoginRequest
	if err := c.BindJSON(&loginRequest); err != nil {
		newErrorResponse(c, http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err))
		return
	}

	if loginRequest.Guid == uuid.Nil {
		newErrorResponse(c, http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest, errors.New("invalid guid param"))
		return
	}
	// fmt.Println(c.RemoteIP())

	tokenPair, err := h.service.Create(c, loginRequest.Guid, c.ClientIP())
	if err != nil {
		newErrorResponse(c, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError, fmt.Errorf("failed to create token pair: %v", err))
		return
	}
	c.JSON(http.StatusOK, tokenPair)
}

func (h *Handler) Refresh(c *gin.Context) {
	var refreshRequest domain.RefreshRequest
	if err := c.BindJSON(&refreshRequest); err != nil {
		newErrorResponse(c, http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err))
		return
	}

	if refreshRequest.SessionID == uuid.Nil {
		newErrorResponse(c, http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest, errors.New("invalid session id param"))
		return
	}

	session, err := h.service.GetSessionById(c, refreshRequest.SessionID)
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			newErrorResponse(c, http.StatusText(http.StatusNotFound),
				http.StatusNotFound, errors.New("session not found"))
			return
		}
		newErrorResponse(c, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError, fmt.Errorf("failed to get session by id: %v", err))
	}

	if session.Ip != c.ClientIP() {
		if err := h.service.NotifyToEmail(session.Ip, c.ClientIP()); err != nil {
			newErrorResponse(c, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError, errors.New("failed to notify to email"))
			return
		}
		newErrorResponse(c, http.StatusText(http.StatusInternalServerError), http.StatusForbidden,
			errors.New("ip address mismatch. Session update is forbidden"))
		return
	}

	refreshedTokenPair, err := h.service.Update(c, refreshRequest.SessionID)
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			newErrorResponse(c, http.StatusText(http.StatusNotFound),
				http.StatusNotFound, errors.New("session not found"))
			return
		}
		newErrorResponse(c, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError, fmt.Errorf("failed to refresh session: %v", err))
		return
	}
	c.JSON(http.StatusOK, refreshedTokenPair)
}

func (h *Handler) GetAllSessions(c *gin.Context) {
	sessions, err := h.service.GetAllSessions(c)
	if err != nil {
		newErrorResponse(c, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError, fmt.Errorf("failed to get sessions list: %v", err))
		return
	}
	if len(sessions) == 0 {
		newErrorResponse(c, http.StatusText(http.StatusNotFound),
			http.StatusNotFound, errors.New("session not found"))
		return
	}
	c.JSON(http.StatusOK, sessions)
}
