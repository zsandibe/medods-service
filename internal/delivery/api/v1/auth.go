package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zsandibe/medods-service/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Login godoc
// @Summary Login
// @Description Creates a new token pairs by guid
// @Tags session
// @Accept  json
// @Produce  json
// @Param   input  body      domain.LoginRequest  false  "Token creation guid"
// @Success 201  {object} domain.TokenPair
// @Failure 400  {object}  errorResponse
// @Failure 500 {object} errorResponse
// @Router /login [post]
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
	c.JSON(http.StatusCreated, tokenPair)
}

// Refresh godoc
// @Summary Refresh token pairs
// @Description Refreshes a new token pairs by taking a session id
// @Tags session
// @Accept  json
// @Produce  json
// @Param   input  body      domain.RefreshRequest  false  "Session updating data"
// @Success 200  {object} domain.TokenPair
// @Failure 400  {object}  errorResponse
// @Faliure 403 {object}  errorResponse
// @Failure 404  {object}  errorResponse
// @Failure 500 {object} errorResponse
// @Router /refresh [put]
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
		newErrorResponse(c, http.StatusText(http.StatusForbidden), http.StatusForbidden,
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

// GetAllSessions godoc
// @Summary Get sessions list
// @Description Getting sessions info
// @Tags session
// @Produce json
// @Success 200 {object} []entity.Session
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /sessions [get]
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
