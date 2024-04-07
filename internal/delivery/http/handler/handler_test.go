package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CodeMaster482/ShortLinkAPI/internal/delivery/http/dto"
	mock_handler "github.com/CodeMaster482/ShortLinkAPI/internal/delivery/http/handler/mocks"
	"github.com/CodeMaster482/ShortLinkAPI/internal/delivery/http/middleware"
	"github.com/CodeMaster482/ShortLinkAPI/internal/model"
	apierror "github.com/CodeMaster482/ShortLinkAPI/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestGetLink(t *testing.T) {
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
		{
			name:           "Not Found Token",
			token:          "token",
			expectedStatus: http.StatusInternalServerError,
			expectedHeader: "",
			expectedBody:   "",
			mockBehaviour: func(usecase *mock_handler.MockLinkUsecase) {
				usecase.EXPECT().GetFullLink(gomock.Any(), "token").Return("", apierror.ErrLinkNotFound).Times(1)
			},
		},
	}

	for _, tc := range testCases {
		test := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gin.SetMode(gin.TestMode)
			usecase := mock_handler.NewMockLinkUsecase(ctrl)
			handler := NewLinkHandler(usecase)

			router := gin.New()
			router.Use(middleware.ErrorMiddleware())
			router.GET("/:key", handler.GetLink)

			test.mockBehaviour(usecase)

			req, err := http.NewRequest(http.MethodGet, "/"+tc.token, http.NoBody)
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
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
		mockBehaviour  func(usecase *mock_handler.MockLinkUsecase)
	}{
		{
			name:           "Valid Token",
			requestBody:    `{"link":"https://example.com"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"short_link":"short","expires_at":"2012-01-10T00:00:00Z"}`,
			mockBehaviour: func(usecase *mock_handler.MockLinkUsecase) {
				usecase.EXPECT().CreateShortLink(
					gomock.Any(),
					&dto.CreateLinkRequest{Link: "https://example.com"},
				).
					Return(&model.Link{
						OriginalLink: "https://example.com",
						ShortLink:    "short",
						Token:        "token",
						ExpiresAt:    time.Date(2012, time.January, 10, 0, 0, 0, 0, time.UTC),
					}, nil).
					Times(1)
			},
		},
		{
			name:           "Creation Error",
			requestBody:    `{"link":"https://example.com"}`,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"internal server error","status":500}`,
			mockBehaviour: func(usecase *mock_handler.MockLinkUsecase) {
				usecase.EXPECT().CreateShortLink(
					gomock.Any(),
					&dto.CreateLinkRequest{Link: "https://example.com"},
				).
					Return(nil, apierror.ErrUnableToCreateLink).
					Times(1)
			},
		},
		{
			name:           "Corruted Request Body",
			requestBody:    `{"link":"https://example.co`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"bad request","status":400}`,
			mockBehaviour:  func(usecase *mock_handler.MockLinkUsecase) {},
		},
		{
			name:           "Empty link value",
			requestBody:    `{"link":""}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"bad request","status":400}`,
			mockBehaviour:  func(usecase *mock_handler.MockLinkUsecase) {},
		},
	}

	for _, tc := range testCases {
		test := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := mock_handler.NewMockLinkUsecase(ctrl)
			handler := NewLinkHandler(usecase)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(middleware.ErrorMiddleware())
			router.POST("/url", handler.CreateLink)

			test.mockBehaviour(usecase)

			req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/url", bytes.NewBufferString(test.requestBody))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d; got %d", tc.expectedStatus, w.Code)
			}

			if tc.expectedBody != "" && w.Body.String() != tc.expectedBody {
				t.Errorf("expected body %q; got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}
