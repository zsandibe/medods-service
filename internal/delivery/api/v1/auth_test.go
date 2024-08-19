package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zsandibe/medods-service/internal/delivery/api/v1/mocks"
	"github.com/zsandibe/medods-service/internal/domain"
)

func setupRouter(handler *Handler) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", handler.Login)
			auth.PUT("/refresh", handler.Refresh)
		}
	}
	return router
}

func TestHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	router := setupRouter(handler)

	tests := []struct {
		name          string
		requestBody   interface{}
		mockExpect    func()
		expectedCode  int
		expectedError string
	}{
		{
			name:          "Invalid JSON",
			requestBody:   "invalid-json",
			mockExpect:    func() {},
			expectedCode:  http.StatusBadRequest,
			expectedError: "invalid request body",
		},
		{
			name:          "Invalid GUID",
			requestBody:   domain.LoginRequest{Guid: uuid.Nil},
			mockExpect:    func() {},
			expectedCode:  http.StatusBadRequest,
			expectedError: "invalid guid param",
		},
		{
			name:        "Internal Server Error",
			requestBody: domain.LoginRequest{Guid: uuid.New()},
			mockExpect: func() {
				mockService.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(domain.TokenPair{}, errors.New("service error")).
					Times(1)
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "failed to create token pair",
		},
		{
			name:        "Success",
			requestBody: domain.LoginRequest{Guid: uuid.New()},
			mockExpect: func() {
				mockService.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(domain.TokenPair{
						AccessToken:  "access-token",
						RefreshToken: []byte("refresh-token"),
					}, nil).
					Times(1)
			},
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(reqBody))
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				var response errorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, response.Code)
				assert.Contains(t, response.Message, tt.expectedError)
			} else {
				var response domain.TokenPair
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, domain.TokenPair{
					AccessToken:  "access-token",
					RefreshToken: []byte("refresh-token"),
				}, response)
			}
		})
	}
}

func TestHandler_Refresh_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	router := setupRouter(handler)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodPut, "/api/v1/auth/refresh", bytes.NewBuffer([]byte("invalid-json")))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Message, "invalid request body")
}

func TestHandler_Refresh_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	sessionID := uuid.New()
	refreshRequest := domain.RefreshRequest{
		SessionID: sessionID,
	}

	mockService.EXPECT().
		GetSessionById(gomock.Any(), sessionID).
		Return(nil, errors.New("service error")).
		Times(1)

	router := setupRouter(handler)
	w := httptest.NewRecorder()
	reqBody, _ := json.Marshal(refreshRequest)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/auth/refresh", bytes.NewBuffer(reqBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Contains(t, response.Message, "failed to get session by id")
}

func TestHandler_Refresh_InvalidSessionID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	router := setupRouter(handler)
	w := httptest.NewRecorder()
	refreshRequest := domain.RefreshRequest{
		SessionID: uuid.Nil,
	}
	reqBody, _ := json.Marshal(refreshRequest)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/auth/refresh", bytes.NewBuffer(reqBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Message, "invalid session id param")
}
