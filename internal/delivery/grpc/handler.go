package grpc

import (
	"context"

	"ShortLinkAPI/internal/delivery/grpc/generated"
	"ShortLinkAPI/internal/delivery/http/dto"
	"ShortLinkAPI/internal/model"
	apierror "ShortLinkAPI/pkg/errors"
)

type LinkUsecase interface {
	GetFullLink(ctx context.Context, token string) (string, error)
	CreateShortLink(ctx context.Context, linkRequest *dto.CreateLinkRequest) (*model.Link, error)
}

type LinkGrpcHandler struct {
	usecase LinkUsecase
	generated.UnimplementedShortLinkServiceServer
}

func NewLinkHandler(usecase LinkUsecase) *LinkGrpcHandler {
	return &LinkGrpcHandler{
		usecase: usecase,
	}
}

func (lgh *LinkGrpcHandler) CreateShortLink(ctx context.Context, request *generated.CreateShortLinkRequest) (*generated.CreateShortLinkResponse, error) {
	addLink := &dto.CreateLinkRequest{Link: request.OriginalLink}
	link, err := lgh.usecase.CreateShortLink(ctx, addLink)
	if err != nil {
		return nil, err
	}

	return &generated.CreateShortLinkResponse{
		ShortLink: link.ShortLink,
		ExpiresAt: link.ExpiresAt.String(),
	}, nil
}

func (lgh *LinkGrpcHandler) GetFullLink(ctx context.Context, request *generated.ShortLinkRequest) (*generated.ShortLinkResponse, error) {
	if request.ShortLink == "" {
		return nil, apierror.BadRequestError()
	}

	link, err := lgh.usecase.GetFullLink(ctx, request.ShortLink)
	if err != nil {
		return nil, err
	}

	return &generated.ShortLinkResponse{
		OriginalLink: link,
	}, nil
}
