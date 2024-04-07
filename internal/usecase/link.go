package usecase

import (
	"ShortLinkAPI/config"
	"ShortLinkAPI/internal/delivery/http/dto"
	"ShortLinkAPI/internal/model"
	apierror "ShortLinkAPI/pkg/errors"
	"context"
	"fmt"
	"net/url"
	"time"
)

type LinkRepository interface {
	GetLink(ctx context.Context, token string) (*model.Link, error)
	GetLinkByOriginal(ctx context.Context, origLink string) (*model.Link, error)
	StoreLink(ctx context.Context, link *model.Link) error
	StartRecalculation(interval time.Duration, deleted chan []string)
	// ShutDown(ctx context.Context) error
}

type Generator interface {
	GenerateShortURL(url string) string
}

type LinkService struct {
	repository      LinkRepository
	generator       Generator
	shortlinkPrefix string
	expiration      time.Time
}

func (service *LinkService) GetFullLink(ctx context.Context, token string) (string, error) {
	link, err := service.repository.GetLink(ctx, token)
	if err != nil {
		return "", err
	}
	return link.OriginalLink, nil
}

func (service *LinkService) CreateShortLink(ctx context.Context, linkRequest *dto.CreateLinkRequest) (*model.Link, error) {
	_, err := url.ParseRequestURI(linkRequest.Link)
	if err != nil {
		return nil, apierror.NewAppError(apierror.ErrUrlNotValid, err)
	}

	link, _ := service.repository.GetLinkByOriginal(ctx, linkRequest.Link)
	if link != nil {
		link.ShortLink = service.shortlinkPrefix + link.Token
		return link, nil
	}

	token := ""
	token = service.generator.GenerateShortURL(linkRequest.Link)

	link = &model.Link{
		OriginalLink: linkRequest.Link,
		Token:        token,
		ExpiresAt:    service.expiration,
		ShortLink:    fmt.Sprintf(service.shortlinkPrefix + token),
	}
	if err := service.repository.StoreLink(ctx, link); err != nil {
		return nil, err
	}

	return link, nil
}

func NewLinkService(cfg *config.Config, repo LinkRepository, strGenerator Generator) *LinkService {
	deleteChan := make(chan []string)
	repo.StartRecalculation(time.Duration(cfg.Service.RecalculationInterval)*time.Hour, deleteChan)
	prefix := fmt.Sprintf("http://%s:%d/url/", cfg.Service.Host, cfg.Service.Port)
	return &LinkService{
		repository:      repo,
		generator:       strGenerator,
		shortlinkPrefix: prefix,
		expiration:      time.Now().Add(time.Duration(24) * time.Hour), //TODO: cfg add
	}
}