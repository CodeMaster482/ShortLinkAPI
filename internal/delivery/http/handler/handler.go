package handler

import (
	"context"
	"net/http"

	"github.com/CodeMaster482/ShortLinkAPI/internal/delivery/http/dto"
	"github.com/CodeMaster482/ShortLinkAPI/internal/model"
	apierror "github.com/CodeMaster482/ShortLinkAPI/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
)

type LinkHandler struct {
	usecase LinkUsecase
}

type LinkUsecase interface {
	GetFullLink(ctx context.Context, token string) (string, error)
	CreateShortLink(ctx context.Context, linkRequest *dto.CreateLinkRequest) (*model.Link, error)
}

func NewLinkHandler(usecase LinkUsecase) *LinkHandler {
	return &LinkHandler{usecase}
}

func (h *LinkHandler) GetLink(ctx *gin.Context) {
	token := ctx.Param("key")

	if token == "" {
		_ = ctx.Error(apierror.BadRequestError())
		return
	}

	link, err := h.usecase.GetFullLink(ctx.Request.Context(), token)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Redirect(http.StatusFound, link)
}

func (h *LinkHandler) CreateLink(ctx *gin.Context) {
	request := &dto.CreateLinkRequest{}
	if err := easyjson.UnmarshalFromReader(ctx.Request.Body, request); err != nil {
		_ = ctx.Error(apierror.BadRequestError())
		return
	}

	if request.Link == "" {
		_ = ctx.Error(apierror.BadRequestError())
		return
	}

	link, err := h.usecase.CreateShortLink(ctx.Request.Context(), request)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := &dto.CreateLinkResponse{
		ShortLink: link.ShortLink,
		ExpiresAt: link.ExpiresAt,
	}

	responseJSON, err := response.MarshalJSON()
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Data(http.StatusOK, "application/json; charset=utf-8", responseJSON)
}
