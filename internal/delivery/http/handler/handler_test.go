package handler

import (
	"ShortLinkAPI/internal/delivery/http/dto"
	mock_handler "ShortLinkAPI/internal/delivery/http/handler/mocks"
	"ShortLinkAPI/internal/model"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"gopkg.in/go-playground/assert.v1"
)

func TestGetLink(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		token          string
		expectedStatus int
		expectedHeader string
		expectedBody   string
		mockBehaviour  func(usecase *mock_handler.MockLinkUsecase)
	}{
		{
			name:           "Valid Token",
			token:          "validToken",
			expectedStatus: http.StatusFound,
			expectedHeader: "https://example.com",
			expectedBody:   "",
			mockBehaviour: func(usecase *mock_handler.MockLinkUsecase) {
				usecase.EXPECT().GetFullLink(gomock.Any(), "validToken").Return("https://example.com", nil).Times(1)
			},
		},
		{
			name:           "Empty Token",
			token:          "",
			expectedStatus: http.StatusNotFound,
			expectedHeader: "",
			expectedBody:   "404 page not found",
			mockBehaviour: func(usecase *mock_handler.MockLinkUsecase) {
				usecase.EXPECT().GetFullLink(gomock.Any(), "").Times(0)
			},
		},
		// Add more test cases as needed
	}

	// Run tests
	for _, tc := range testCases {
		test := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			ctrl := gomock.NewController(t)

			gin.SetMode(gin.TestMode)
			usecase := mock_handler.NewMockLinkUsecase(ctrl)
			handler := NewLinkHandler(usecase)
			router := gin.New()
			router.GET("/:key", handler.GetLink)

			test.mockBehaviour(usecase)

			req, err := http.NewRequest(http.MethodGet, "/"+tc.token, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d; got %d", tc.expectedStatus, w.Code)
			}

			if tc.expectedHeader != "" && w.Header().Get("Location") != tc.expectedHeader {
				t.Errorf("expected header %q; got %q", tc.expectedHeader, w.Header().Get("Location"))
			}

			if tc.expectedBody != "" && w.Body.String() != tc.expectedBody {
				t.Errorf("expected body %q; got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestCreateLink(t *testing.T) {
	testCases := []struct {
		name             string
		requestBody      string
		expectedStatus   int
		expectedResponse dto.CreateLinkResponse
		mockBehaviour    func(usecase *mock_handler.MockLinkUsecase)
	}{
		{
			name:             "Valid Request",
			requestBody:      `{"Link": "https://example.com"}`,
			expectedStatus:   http.StatusOK,
			expectedResponse: dto.CreateLinkResponse{ShortLink: "shortLink", ExpiresAt: time.Now().Add(24 * time.Hour)},
			mockBehaviour: func(usecase *mock_handler.MockLinkUsecase) {
				usecase.EXPECT().CreateShortLink(gomock.Any(), gomock.Any()).Return(&model.Link{OriginalLink: "https://example.com", ShortLink: "shortLink", ExpiresAt: time.Now().Add(24 * time.Hour)}, nil).Times(1)
			},
		},
	}

	// Run tests
	for _, tc := range testCases {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gin.SetMode(gin.TestMode)
			r := gin.New()

			mockUsecase := mock_handler.NewMockLinkUsecase(ctrl)
			handler := NewLinkHandler(mockUsecase)

			r.POST("/url", handler.CreateLink)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/url", bytes.NewBufferString(test.requestBody))
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			test.mockBehaviour(mockUsecase)

			// Call the function we're testing
			handler.CreateLink(ctx)

			// Check the response
			assert.Equal(t, test.expectedStatus, w.Code)

			if test.expectedStatus == http.StatusOK {
				var response dto.CreateLinkResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("could not unmarshal response: %v", err)
				}

				assert.Equal(t, test.expectedResponse, response)
			}
		})
	}
}
