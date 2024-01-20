package driver

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type RDBDriver interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	Beginx() (*sqlx.Tx, error)
	Rollback() error
}

type txDriver struct {
	*sqlx.Tx
}

func NewTxDriver(tx *sqlx.Tx) RDBDriver {
	return txDriver{tx}
}

// dummy
func (driver txDriver) Beginx() (*sqlx.Tx, error) {
	return driver.Tx, nil
}

type rdbDriver struct {
	*sqlx.DB
}

// dummy
func (driver rdbDriver) Rollback() error {
	return nil
}

func NewRDBDriver() (RDBDriver, error) {
	db, err := sqlx.Open(
		"mysql",
		fmt.Sprintf("root:root@tcp(127.0.0.1:3306)/super_shiharai_kun?charset=utf8&parseTime=true"), // TODO: 接続先は変更できるようにする
	)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &rdbDriver{db}, nil
}
