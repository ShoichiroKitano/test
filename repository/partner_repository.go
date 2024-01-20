package repository

import (
	"super-shiharai-kun/entity"
)

// TODO: DBに保存して使えるようにする
func (repo *RDBRepository) FindPartnerByCorporateID(corporateID uint64) (entity.Partners, error) {
	return []*entity.Partner{
		{ID: 1000, CorporateID: 100},
		{ID: 1001, CorporateID: 100},
	}, nil
}
