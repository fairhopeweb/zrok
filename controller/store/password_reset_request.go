package store

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type PasswordResetRequest struct {
	Model
	Token     string
	AccountId int
}

func (self *Store) CreatePasswordResetRequest(prr *PasswordResetRequest, tx *sqlx.Tx) (int, error) {
	stmt, err := tx.Prepare("insert into password_reset_requests (account_id, token) values ($1, $2) ON CONFLICT(account_id) DO UPDATE SET token=$2 returning id")
	if err != nil {
		return 0, errors.Wrap(err, "error preparing password_reset_requests insert statement")
	}
	var id int
	if err := stmt.QueryRow(prr.AccountId, prr.Token).Scan(&id); err != nil {
		return 0, errors.Wrap(err, "error executing password_reset_requests insert statement")
	}
	return id, nil
}

func (self *Store) FindPasswordResetRequestWithToken(token string, tx *sqlx.Tx) (*PasswordResetRequest, error) {
	prr := &PasswordResetRequest{}
	if err := tx.QueryRowx("select * from password_reset_requests where token = $1", token).StructScan(prr); err != nil {
		return nil, errors.Wrap(err, "error selecting password_reset_requests by token")
	}
	return prr, nil
}

func (self *Store) FindExpiredPasswordResetRequests(before time.Time, limit int, tx *sqlx.Tx) ([]*PasswordResetRequest, error) {
	var sql string
	switch self.cfg.Type {
	case "postgres":
		sql = "select * from password_reset_requests where created_at < $1 limit %d for update"

	case "sqlite3":
		sql = "select * from password_reset_requests where created_at < $1 limit %d"
	default:
		return nil, errors.Errorf("unknown database type '%v'", self.cfg.Type)
	}

	rows, err := tx.Queryx(fmt.Sprintf(sql, limit), before)
	if err != nil {
		return nil, errors.Wrap(err, "error selecting expired password_reset_requests")
	}
	var prrs []*PasswordResetRequest
	for rows.Next() {
		prr := &PasswordResetRequest{}
		if err := rows.StructScan(prr); err != nil {
			return nil, errors.Wrap(err, "error scanning password_reset_request")
		}
		prrs = append(prrs, prr)
	}
	return prrs, nil
}

func (self *Store) DeletePasswordResetRequest(id int, tx *sqlx.Tx) error {
	stmt, err := tx.Prepare("delete from password_reset_requests where id = $1")
	if err != nil {
		return errors.Wrap(err, "error preparing password_reset_requests delete statement")
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return errors.Wrap(err, "error executing password_reset_requests delete statement")
	}
	return nil
}

func (self *Store) DeleteMultiplePasswordResetRequests(ids []int, tx *sqlx.Tx) error {
	if len(ids) == 0 {
		return nil
	}

	anyIds := make([]any, len(ids))
	indexes := make([]string, len(ids))

	for i, id := range ids {
		anyIds[i] = id
		indexes[i] = fmt.Sprintf("$%d", i+1)
	}

	stmt, err := tx.Prepare(fmt.Sprintf("delete from password_reset_requests where id in (%s)", strings.Join(indexes, ",")))
	if err != nil {
		return errors.Wrap(err, "error preparing password_reset_requests delete multiple statement")
	}
	_, err = stmt.Exec(anyIds...)
	if err != nil {
		return errors.Wrap(err, "error executing password_reset_requests delete multiple statement")
	}
	return nil
}
