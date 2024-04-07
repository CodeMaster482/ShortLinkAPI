package redis

import (
	"ShortLinkAPI/internal/model"
	apierror "ShortLinkAPI/pkg/errors"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestStoreLink(t *testing.T) {
	mockClient, mock := redismock.NewClientMock()

	repo := &LinkRedisStorage{
		Client: mockClient,
	}

	token := "short"
	longUrl := "https://www.example.com"
	expirationTime := time.Now().Add(1 * time.Hour)

	mock.ExpectSet(token, longUrl, 0).SetVal(token)
	mock.ExpectExpire(token, expirationTime.Sub(time.Now())).SetVal(true)

	err := repo.StoreLink(
		context.TODO(),
		&model.Link{
			OriginalLink: longUrl,
			Token:        token,
			ExpiresAt:    expirationTime,
		},
	)

	assert.Nil(t, err, "Expected no error, got %v", err)
	assert.NoError(t, mock.ExpectationsWereMet(), "Expectations were not met")
}

func TestSaveUrl_SetError(t *testing.T) {
	mockClient, mock := redismock.NewClientMock()

	repo := &LinkRedisStorage{
		Client: mockClient,
	}

	token := "short"
	longUrl := "https://www.example.com"
	expirationTime := time.Now().Add(1 * time.Hour)

	expectedError := fmt.Errorf("set error")
	mock.ExpectSet(token, longUrl, 0).SetErr(expectedError)

	err := repo.StoreLink(
		context.TODO(),
		&model.Link{
			OriginalLink: longUrl,
			Token:        token,
			ExpiresAt:    expirationTime,
		},
	)

	assert.EqualError(t, err, expectedError.Error(), "Expected error does not match actual error")
}

func TestGetLink_Success(t *testing.T) {
	mockClient, mock := redismock.NewClientMock()

	repo := &LinkRedisStorage{
		Client: mockClient,
	}

	url := "short"
	longUrl := "https://www.example.com"
	mock.ExpectGet(url).SetVal(longUrl)

	result, err := repo.GetLink(context.TODO(), url)

	assert.Nil(t, err, "Expected no error, got %v", err)
	assert.Equal(t, longUrl, result, "Expected long URL %s, got %s", longUrl, result)

	assert.NoError(t, mock.ExpectationsWereMet(), "Expectations were not met")
}

func TestGetLink_NotFound(t *testing.T) {
	mockClient, mock := redismock.NewClientMock()

	repo := &LinkRedisStorage{
		Client: mockClient,
	}

	token := "short"
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
	mockClient, mock := redismock.NewClientMock()

	repo := &LinkRedisStorage{
		Client: mockClient,
	}

	url := "short"
	expectedError := errors.New("something went wrong")
	mock.ExpectGet(url).SetErr(expectedError)

	result, err := repo.GetLink(context.TODO(), url)

	assert.Error(t, err, "Expected an error")
	assert.Equal(t, result, "", "Expected %s, got %s", result, "")
	assert.Equal(t, expectedError, err, "Expected %v, got %v", expectedError, err)

	assert.NoError(t, mock.ExpectationsWereMet(), "Expectations were not met")
}
