package repository

import (
	"super-shiharai-kun/entity"
)

// TODO: DBに保存して使えるようにする
func (repo *RDBRepository) FindCorporateByUserID(userID uint64) (*entity.Corporate, error) {
	return &entity.Corporate{ID: 100}, nil
}
