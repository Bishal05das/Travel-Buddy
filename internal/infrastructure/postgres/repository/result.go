package repository

import (
	"database/sql"
	"errors"
)

// requireAffected turns an UPDATE/DELETE that matched no rows into notFoundMsg.
func requireAffected(res sql.Result, notFoundMsg string) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New(notFoundMsg)
	}
	return nil
}
