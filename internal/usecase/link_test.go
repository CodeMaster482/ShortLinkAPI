package usecase

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/CodeMaster482/ShortLinkAPI/internal/delivery/http/dto"
	"github.com/CodeMaster482/ShortLinkAPI/internal/model"
	mock_usecase "github.com/CodeMaster482/ShortLinkAPI/internal/usecase/mocks"
	apierror "github.com/CodeMaster482/ShortLinkAPI/pkg/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	prefix = "http://localhost:8080/url/"
)

func TestLinkService_GetOriginalLink(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		expectedLink  string
		expectedError error
		token         string

		mockBehaviour func(repository *mock_usecase.MockLinkRepository,
			generator *mock_usecase.MockGenerator,
			token, link string)
	}{
		{
			name:         "Success",
			expectedLink: "http://wikipedia.org",
			token:        "qwerty123_",
			mockBehaviour: func(repository *mock_usecase.MockLinkRepository,
				generator *mock_usecase.MockGenerator,
				token, link string) {
				linkReturned := &model.Link{
					OriginalLink: link,
				}
				repository.EXPECT().GetLink(gomock.Any(), gomock.Any()).Return(linkReturned, nil)
			},
		},
		{
			name:          "Not found",
			expectedLink:  "",
			expectedError: apierror.NotFoundError(),
			mockBehaviour: func(repository *mock_usecase.MockLinkRepository,
				generator *mock_usecase.MockGenerator,
				token, link string) {
				repository.EXPECT().GetLink(gomock.Any(), gomock.Any()).Return(nil, apierror.NotFoundError())
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_usecase.NewMockLinkRepository(ctrl)
			mockGenerator := mock_usecase.NewMockGenerator(ctrl)

			test.mockBehaviour(mockRepo, mockGenerator, test.token, test.expectedLink)

			usecase := LinkService{
				repository:      mockRepo,
				generator:       mockGenerator,
				shortlinkPrefix: prefix,
			}

			link, err := usecase.GetFullLink(context.TODO(), test.token)
			if test.expectedError != nil {
				require.ErrorAs(t, err, &test.expectedError)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.expectedLink, link)
		})
	}
}

func TestLinkService_CreateShortLink(t *testing.T) {
	tests := []struct {
		name          string
		expectedLink  *model.Link
		expectedError error
		token         string
		dto           *dto.CreateLinkRequest

		mockBehaviour func(repository *mock_usecase.MockLinkRepository,
			generator *mock_usecase.MockGenerator,
			dto *dto.CreateLinkRequest, link *model.Link)
	}{
		{
			name: "Success",
			dto: &dto.CreateLinkRequest{
				Link: "http://wikipedia.org",
			},
			expectedLink: &model.Link{
				OriginalLink: "http://wikipedia.org",
				Token:        "qwerty123_",
				ExpiresAt:    time.Now(),
				ShortLink:    prefix + "qwerty123_",
			},
			mockBehaviour: func(repository *mock_usecase.MockLinkRepository, generator *mock_usecase.MockGenerator, dto *dto.CreateLinkRequest, link *model.Link) {
				repository.EXPECT().GetLinkByOriginal(gomock.Any(), dto.Link).Return(nil, nil)
				generator.EXPECT().GenerateShortURL(gomock.Any()).Return(link.Token).AnyTimes()
				repository.EXPECT().StoreLink(gomock.Any(), gomock.Any()).Return(nil)
			},
		}, {
			name: "invalid uri",
			dto: &dto.CreateLinkRequest{
				Link: "bag",
			},
			expectedLink:  nil,
			expectedError: apierror.NewAPIError(apierror.ErrURLNotValid, fmt.Errorf("empty URL")),
			mockBehaviour: func(repository *mock_usecase.MockLinkRepository, generator *mock_usecase.MockGenerator, dto *dto.CreateLinkRequest, link *model.Link) {
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_usecase.NewMockLinkRepository(ctrl)
			mockGenerator := mock_usecase.NewMockGenerator(ctrl)

			test.mockBehaviour(mockRepo, mockGenerator, test.dto, test.expectedLink)

			usecase := LinkService{
				repository:      mockRepo,
				generator:       mockGenerator,
				shortlinkPrefix: prefix,
			}

			link, err := usecase.CreateShortLink(context.TODO(), test.dto)
			if test.expectedError != nil {
				require.ErrorAs(t, err, &test.expectedError)
				return
			}

			require.NoError(t, err)

			if link != nil {
				link.ExpiresAt = time.Now()
			}

			require.Equal(t, test.expectedLink.ShortLink, link.ShortLink)
			require.Equal(t, test.expectedLink.ShortLink, link.ShortLink)
		})
	}
}
