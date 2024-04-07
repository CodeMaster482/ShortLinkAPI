package grpc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CodeMaster482/ShortLinkAPI/internal/delivery/grpc"
	"github.com/CodeMaster482/ShortLinkAPI/internal/delivery/grpc/generated"
	"github.com/CodeMaster482/ShortLinkAPI/internal/delivery/http/dto"
	"github.com/CodeMaster482/ShortLinkAPI/internal/model"

	mock_handler "github.com/CodeMaster482/ShortLinkAPI/internal/delivery/http/handler/mocks"
	apierror "github.com/CodeMaster482/ShortLinkAPI/pkg/errors"

	"github.com/golang/mock/gomock"
)

func TestCreateShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_handler.NewMockLinkUsecase(ctrl)
	handler := grpc.NewLinkHandler(mockUsecase)

	ctx := context.Background()
	request := &generated.CreateShortLinkRequest{
		OriginalLink: "http://example.com",
	}

	expectedShortLink := "http://short.link/abc123"
	expectedExpiresAt := time.Now().Add(24 * time.Hour)

	mockUsecase.EXPECT().
		CreateShortLink(ctx, &dto.CreateLinkRequest{Link: request.OriginalLink}).
		Return(&model.Link{ShortLink: expectedShortLink, ExpiresAt: expectedExpiresAt}, nil)

	response, err := handler.CreateShortLink(ctx, request)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expectedResponse := &generated.CreateShortLinkResponse{
		ShortLink: expectedShortLink,
		ExpiresAt: expectedExpiresAt.String(),
	}

	if response.ShortLink != expectedResponse.ShortLink || response.ExpiresAt != expectedResponse.ExpiresAt {
		t.Errorf("Unexpected response. Expected: %v, Got: %v", expectedResponse, response)
	}
}

func TestGetFullLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_handler.NewMockLinkUsecase(ctrl)
	handler := grpc.NewLinkHandler(mockUsecase)

	ctx := context.Background()
	request := &generated.ShortLinkRequest{
		ShortLink: "http://short.link/abc123",
	}

	expectedOriginalLink := "http://example.com"

	mockUsecase.EXPECT().
		GetFullLink(ctx, request.ShortLink).
		Return(expectedOriginalLink, nil)

	response, err := handler.GetFullLink(ctx, request)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expectedResponse := &generated.ShortLinkResponse{
		OriginalLink: expectedOriginalLink,
	}

	if response.OriginalLink != expectedResponse.OriginalLink {
		t.Errorf("Unexpected response. Expected: %v, Got: %v", expectedResponse, response)
	}
}

func TestCreateShortLink_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_handler.NewMockLinkUsecase(ctrl)
	handler := grpc.NewLinkHandler(mockUsecase)

	ctx := context.Background()
	request := &generated.CreateShortLinkRequest{
		OriginalLink: "http://example.com",
	}

	expectedError := errors.New("mock error")

	mockUsecase.EXPECT().
		CreateShortLink(ctx, &dto.CreateLinkRequest{Link: request.OriginalLink}).
		Return(nil, expectedError)

	_, err := handler.CreateShortLink(ctx, request)
	if err == nil {
		t.Error("Expected an error but got nil")
		return
	}

	if err != expectedError {
		t.Errorf("Unexpected error. Expected: %v, Got: %v", expectedError, err)
	}
}

func TestGetFullLink_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_handler.NewMockLinkUsecase(ctrl)
	handler := grpc.NewLinkHandler(mockUsecase)

	ctx := context.Background()
	request := &generated.ShortLinkRequest{
		ShortLink: "http://short.link/abc123",
	}

	expectedError := apierror.BadRequestError()

	mockUsecase.EXPECT().
		GetFullLink(ctx, request.ShortLink).
		Return("", expectedError)

	_, err := handler.GetFullLink(ctx, request)
	if err == nil {
		t.Error("Expected an error but got nil")
		return
	}

	if err != expectedError {
		t.Errorf("Unexpected error. Expected: %v, Got: %v", expectedError, err)
	}
}
