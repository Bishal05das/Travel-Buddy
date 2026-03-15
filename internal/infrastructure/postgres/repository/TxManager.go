package repository

import (
	"context"
	"errors"

	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/jmoiron/sqlx"
)

type txKey struct{}

type TxManager struct {
	DB *sqlx.DB
}

func NewTxManager(DB *sqlx.DB) port.TxManager {
	return &TxManager{
		DB: DB,
	}
}

func (m *TxManager) WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx, err := m.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.New("tx rollback failed: " + rbErr.Error() + ", original error: " + err.Error())
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func GetTx(ctx context.Context) (*sqlx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sqlx.Tx)
	return tx, ok
}
