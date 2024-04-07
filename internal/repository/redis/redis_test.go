package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/CodeMaster482/ShortLinkAPI/internal/model"
	apierror "github.com/CodeMaster482/ShortLinkAPI/pkg/errors"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

const (
	testURL   = "https://www.example.com"
	testToken = "short"
)

func TestStoreLink(t *testing.T) {
	t.Parallel()

	mockClient, mock := redismock.NewClientMock()
	repo := &LinkRedisStorage{
		Client: mockClient,
	}

	token := testToken
	longURL := testURL
	expirationTime := time.Now().Add(1 * time.Hour)

	mock.ExpectSet(token, longURL, 0).SetVal(token)
	mock.ExpectExpire(token, time.Until(expirationTime)).SetVal(true)

	err := repo.StoreLink(
		context.TODO(),
		&model.Link{
			OriginalLink: longURL,
			Token:        token,
			ExpiresAt:    expirationTime,
		},
	)

	assert.Nil(t, err, "Expected no error, got %v", err)
	assert.NoError(t, mock.ExpectationsWereMet(), "Expectations were not met")
}

func TestSaveLink_SetError(t *testing.T) {
	t.Parallel()
	mockClient, mock := redismock.NewClientMock()

	repo := &LinkRedisStorage{
		Client: mockClient,
	}

	token := testToken
	longURL := testURL
	expirationTime := time.Now().Add(1 * time.Hour)

	expectedError := fmt.Errorf("set error")
	mock.ExpectSet(token, longURL, 0).SetErr(expectedError)

	err := repo.StoreLink(
		context.TODO(),
		&model.Link{
			OriginalLink: longURL,
			Token:        token,
			ExpiresAt:    expirationTime,
		},
	)

	assert.EqualError(t, err, expectedError.Error(), "Expected error does not match actual error")
}

func TestGetLink_Success(t *testing.T) {
	t.Parallel()
	mockClient, mock := redismock.NewClientMock()

	repo := &LinkRedisStorage{
		Client: mockClient,
	}

	url := testToken
	longURL := testURL
	mock.ExpectGet(url).SetVal(longURL)

	result, err := repo.GetLink(context.TODO(), url)

	assert.Nil(t, err, "Expected no error, got %v", err)
	assert.Equal(t, longURL, result, "Expected long URL %s, got %s", longURL, result)

	assert.NoError(t, mock.ExpectationsWereMet(), "Expectations were not met")
}

func TestGetLink_NotFound(t *testing.T) {
	t.Parallel()
	mockClient, mock := redismock.NewClientMock()

	repo := &LinkRedisStorage{
		Client: mockClient,
	}

	token := testToken
	expectedError := redis.Nil
	mock.ExpectGet(token).SetErr(expectedError)

	result, err := repo.GetLink(context.TODO(), token)

	assert.Error(t, err, "Expected an error")
	assert.Equal(t, result, "", "Expected %s, got %s", result, "")

	assert.IsType(t, apierror.ErrLinkNotFound, err, "Expected error type to be NoSuchLink")
	assert.Equal(t, apierror.ErrLinkNotFound.Error(), err.Error(), "Expected error message %q, got %q", fmt.Sprintf("No such url link: %v", expectedError), err.Error())

	assert.NoError(t, mock.ExpectationsWereMet(), "Expectations were not met")
}

func TestGetUrl_Error(t *testing.T) {
	t.Parallel()
	mockClient, mock := redismock.NewClientMock()

	repo := &LinkRedisStorage{
		Client: mockClient,
	}

	url := testToken
	expectedError := fmt.Errorf("something went wrong")
	mock.ExpectGet(url).SetErr(expectedError)

	result, err := repo.GetLink(context.TODO(), url)

	assert.Error(t, err, "Expected an error")
	assert.Equal(t, result, "", "Expected %s, got %s", result, "")
	assert.Equal(t, expectedError, err, "Expected %v, got %v", expectedError, err)

	assert.NoError(t, mock.ExpectationsWereMet(), "Expectations were not met")
}
