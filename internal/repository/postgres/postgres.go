package postgres

import (
	"context"
	"errors"
	"time"

	"ShortLinkAPI/internal/model"
	"ShortLinkAPI/internal/utils"
	apierror "ShortLinkAPI/pkg/errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBConn interface {
	// Conn() *pgx.Conn
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) error
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type LinkStorage struct {
	db DBConn
}

func (store *LinkStorage) GetLink(ctx context.Context, token string) (*model.Link, error) {
	query := `SELECT s.original_link, s.token, s.expires_at FROM link s WHERE s.token = $1;`
	link := &model.Link{}
	expireStr := ""

	err := store.db.QueryRow(context.Background(), query, token).Scan(&link.OriginalLink, &link.Token, &expireStr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.ErrLinkNotFound
		}
		return nil, err
	}

	link.ExpiresAt, _ = time.Parse(time.RFC3339Nano, expireStr)

	return link, nil
}

func (store *LinkStorage) GetLinkByOriginal(ctx context.Context, origLink string) (*model.Link, error) {
	query := `SELECT s.original_link, s.token, s.expires_at FROM link s WHERE s.original_link = $1;`
	link := &model.Link{}

	err := store.db.QueryRow(context.Background(), query, origLink).Scan(&link.OriginalLink, &link.Token, &link.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NotFoundError()
		}
		return nil, apierror.InternalError(err)
	}

	return link, nil
}

func (store *LinkStorage) StoreLink(ctx context.Context, link *model.Link) error {
	query := `INSERT INTO link (original, short, expiration_time) VALUES ($1, $2, $3);`

	_, err := store.db.Exec(context.Background(), query, link.OriginalLink, link.ShortLink, link.ExpiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (store *LinkStorage) StartRecalculation(interval time.Duration, deleted chan []string) {
	query := `DELETE FROM link WHERE expires_at < $1 RETURNING token`
	ticker := time.NewTicker(interval)
	go func() {
		for {
			<-ticker.C
			rows, err := store.db.Query(context.Background(), query, utils.CurrentTimeString())
			if err != nil {
				continue
			}
			var del []string
			for rows.Next() {
				var deletedToken string
				err := rows.Scan(&deletedToken)
				if err != nil {
					continue
				}
				del = append(del, deletedToken)
			}
			deleted <- del
		}
	}()
}

// func (store *LinkStorage) ShutDown(ctx context.Context) error {
// 	return store.db.Conn().Close(ctx)
// }

func NewLinkStorage(db DBConn) *LinkStorage {
	return &LinkStorage{db}
}