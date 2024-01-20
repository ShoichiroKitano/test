package repository

import (
	"super-shiharai-kun/driver"
)

type RDBRepository struct {
	Driver driver.RDBDriver
}

func (repo *RDBRepository) Tx(operation func(tx RDBRepository) error) error {
	tx, err := repo.Driver.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()
	err = operation(RDBRepository{Driver: driver.NewTxDriver(tx)})
	if err != nil {
		if rollbackErr := repo.Driver.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	return tx.Commit()
}
