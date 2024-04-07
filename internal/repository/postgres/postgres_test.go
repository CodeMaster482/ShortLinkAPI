package postgres

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	apierror "github.com/CodeMaster482/ShortLinkAPI/pkg/errors"

	"github.com/CodeMaster482/ShortLinkAPI/internal/model"
	"github.com/CodeMaster482/ShortLinkAPI/internal/utils"

	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

const (
	getLinkByToken    = `SELECT s.original_link, s.token, s.expires_at FROM link s WHERE s.token = $1;`
	getLinkByFullLink = `SELECT s.original_link, s.token, s.expires_at FROM link s WHERE s.original_link = $1;`
	addLink           = `INSERT INTO link (original, short, expiration_time) VALUES ($1, $2, $3);`
)

func TestPostgreSQLRepository_StoreLink(t *testing.T) {
	timeLink := time.Now().Add(24 * time.Hour)
	testCases := []struct {
		name           string
		longUrl        string
		shortUrl       string
		expirationTime time.Time
		link           model.Link
		expectQuery    string
		expectArgs     []interface{}
		expectError    error
	}{
		{
			name: "Valid case",
			link: model.Link{
				OriginalLink: "http://example.com",
				Token:        "abc123",
				ExpiresAt:    timeLink,
			},
			expectQuery: addLink,
			expectError: nil,
		},
		{
			name: "Error case",
			link: model.Link{
				OriginalLink: "http://example.com",
				Token:        "abc123",
				ExpiresAt:    timeLink,
			},
			expectQuery: addLink,
			expectError: errors.New("mock error"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mock, _ := pgxmock.NewPool()

			repo := &LinkStorage{
				db: mock,
			}

			escapedQuery := regexp.QuoteMeta(tc.expectQuery)

			mock.ExpectExec(escapedQuery).
				WithArgs(tc.longUrl, tc.shortUrl, tc.expirationTime).
				WillReturnResult(pgxmock.NewResult("INSERT", 1)).
				WillReturnError(tc.expectError)

			err := repo.StoreLink(
				context.TODO(),
				&model.Link{
					OriginalLink: tc.longUrl,
					Token:        tc.shortUrl,
					ExpiresAt:    tc.expirationTime,
				},
			)

			if tc.expectError != nil {
				assert.EqualError(t, err, tc.expectError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetLinkByOriginal(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		origLink     string
		rows         *pgxmock.Rows
		expectError  error
		errorPgx     error
		expectedLink *model.Link
	}{
		{
			name:     "Link exists",
			origLink: "http://example.com",
			rows: pgxmock.NewRows([]string{"original_link", "token", "expires_at"}).
				AddRow("http://example.com", "abc123", time.Date(2012, time.January, 10, 0, 0, 0, 0, time.UTC)),
			expectError: nil,
			expectedLink: &model.Link{
				OriginalLink: "http://example.com",
				Token:        "abc123",
				ExpiresAt:    time.Date(2012, time.January, 10, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:         "Link does not exist",
			origLink:     "http://nonexistent.com",
			rows:         pgxmock.NewRows([]string{}),
			errorPgx:     pgx.ErrNoRows,
			expectError:  apierror.ErrLinkNotFound,
			expectedLink: nil,
		},
		{
			name:         "Internal error",
			origLink:     "http://example.com",
			rows:         pgxmock.NewRows([]string{}),
			errorPgx:     errors.New("mock error"),
			expectError:  apierror.ErrInternalServer,
			expectedLink: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mock, _ := pgxmock.NewPool()
			repo := &LinkStorage{
				db: mock,
			}

			escapedQuery := regexp.QuoteMeta(getLinkByFullLink)
			mock.ExpectQuery(escapedQuery).
				WithArgs(tc.origLink).
				WillReturnRows(tc.rows).
				WillReturnError(tc.errorPgx)

			result, err := repo.GetLinkByOriginal(context.Background(), tc.origLink)

			if !errors.Is(err, tc.expectError) {
				t.Errorf("unexpected error, expected: %v, got: %v", tc.expectError, err)
			}

			if !reflect.DeepEqual(result, tc.expectedLink) {
				t.Errorf("unexpected result, expected: %v, got: %v", tc.expectedLink, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestLinkStorage_GetLink(t *testing.T) {
	testCases := []struct {
		name        string
		token       string
		rows        *pgxmock.Rows
		expectError error
		errorPgx    error
		result      *model.Link
	}{
		{
			name:  "Valid case",
			token: "abc123",
			rows: pgxmock.NewRows([]string{"original_link", "token", "expires_at"}).
				AddRow("www.youtube.com", "short", "2012-01-10T00:00:00Z"),
			expectError: nil,
			result: &model.Link{
				OriginalLink: "www.youtube.com",
				Token:        "short",
				ExpiresAt:    time.Date(2012, time.January, 10, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:        "No such link",
			token:       "nonexistent",
			rows:        pgxmock.NewRows([]string{}),
			errorPgx:    pgx.ErrNoRows,
			expectError: apierror.ErrLinkNotFound,
		},
		{
			name:        "Error case",
			token:       "abc123",
			rows:        pgxmock.NewRows([]string{}),
			expectError: errors.New("mock error"),
			errorPgx:    errors.New("mock error"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mock, mockErr := pgxmock.NewPool()

			if mockErr != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", mockErr)
			}

			repo := &LinkStorage{
				db: mock,
			}

			escapedQuery := regexp.QuoteMeta("SELECT s.original_link, s.token, s.expires_at FROM link s WHERE s.token = $1")

			mock.ExpectQuery(escapedQuery).
				WithArgs(tc.token).
				WillReturnRows(tc.rows).
				WillReturnError(tc.errorPgx)

			result, err := repo.GetLink(context.Background(), tc.token)

			if tc.expectError != nil {
				assert.EqualError(t, err, tc.expectError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.result, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestStartRecalculation(t *testing.T) {
	t.Parallel()

	mock, mockErr := pgxmock.NewPool()
	if mockErr != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", mockErr)
	}

	// Create LinkStorage with mocked database
	store := &LinkStorage{
		db: mock,
	}

	// Set up mock expectations
	query := regexp.QuoteMeta("DELETE FROM link WHERE expires_at < $1 RETURNING token")
	mock.ExpectQuery(query).
		WithArgs(utils.CurrentTimeString()).
		WillReturnRows(pgxmock.NewRows([]string{"token"}).AddRow("deleted_token"))

	// Set up channel for receiving deleted tokens
	deleted := make(chan []string)
	defer close(deleted)

	// Start recalculation goroutine
	go store.StartRecalculation(time.Second, deleted)

	// Wait for some time to allow the goroutine to execute
	time.Sleep(5 * time.Second)

	// Check if the expected token is received on the channel
	select {
	case deletedTokens := <-deleted:
		// Check if the expected token is received
		if len(deletedTokens) != 1 || deletedTokens[0] != "deleted_token" {
			t.Errorf("unexpected deleted tokens: %v", deletedTokens)
		}
	default:
		t.Error("no token received on the channel")
	}

	// Ensure all expectations are met
	assert.NoError(t, mock.ExpectationsWereMet(), "not all SQL expectations were met")
}

// func TestPostgreSQLRepository_UrlExistsShort(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		shortUrl       string
// 		rows           *pgxmock.Rows
// 		expectError    error
// 		expectedResult string
// 		errorPgx       error
// 	}{
// 		{
// 			name:           "URL exists",
// 			shortUrl:       "existing",
// 			rows:           pgxmock.NewRows([]string{"original"}).AddRow("www.example.com"),
// 			expectError:    nil,
// 			expectedResult: "www.example.com",
// 		},
// 		{
// 			name:           "No such URL",
// 			shortUrl:       "nonexistent",
// 			rows:           pgxmock.NewRows([]string{}),
// 			expectError:    nil,
// 			expectedResult: "",
// 			errorPgx:       pgx.ErrNoRows,
// 		},
// 		{
// 			name:           "Error case",
// 			shortUrl:       "error",
// 			rows:           pgxmock.NewRows([]string{}),
// 			expectError:    errors.New("mock error"),
// 			errorPgx:       errors.New("mock error"),
// 			expectedResult: "",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mock, _ := pgxmock.NewPool()

// 			repo := &PostgreSQLRepository{
// 				Pool: mock,
// 			}

// 			escapedQuery := regexp.QuoteMeta(getLongUrl)

// 			mock.ExpectQuery(escapedQuery).
// 				WithArgs(tt.shortUrl).
// 				WillReturnRows(tt.rows).
// 				WillReturnError(tt.errorPgx)

// 			result, err := repo.UrlExistsShort(tt.shortUrl)

// 			if tt.expectError != nil {
// 				assert.EqualError(t, err, tt.expectError.Error())
// 				assert.Equal(t, tt.expectedResult, result)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.expectedResult, result)
// 			}

// 			assert.NoError(t, mock.ExpectationsWereMet())
// 		})
// 	}
// }

// func TestPostgreSQLRepository_Clear(t *testing.T) {
// 	tests := []struct {
// 		name             string
// 		rowsAffected     int64
// 		expectedErrorMsg error
// 		errorPgx         error
// 	}{
// 		{
// 			name:         "Success case",
// 			rowsAffected: 5,
// 			errorPgx:     nil,
// 		},
// 		{
// 			name:             "Error case",
// 			rowsAffected:     0,
// 			errorPgx:         errors.New("mock error"),
// 			expectedErrorMsg: fmt.Errorf("error deleting outdated records: %v", errors.New("mock error")),
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mock, _ := pgxmock.NewPool()

// 			log, _ := logger.NewLogger("test.log")

// 			repo := &PostgreSQLRepository{
// 				Pool: mock,
// 				Log:  *log,
// 			}

// 			mock.ExpectExec(regexp.QuoteMeta(deleteOldLink)).
// 				WillReturnResult(pgxmock.NewResult("DELETE", tt.rowsAffected)).
// 				WillReturnError(tt.errorPgx)

// 			err := repo.Clear()

// 			if tt.errorPgx != nil {
// 				assert.Error(t, err)
// 				assert.EqualError(t, err, tt.expectedErrorMsg.Error())
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			assert.NoError(t, mock.ExpectationsWereMet())
// 		})
// 	}
// }
