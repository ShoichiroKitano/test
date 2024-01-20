package repository

import (
	"github.com/shopspring/decimal"
	"super-shiharai-kun/entity"
	"time"
)

// TODO: DBに保存して使えるようにする
func (repo *RDBRepository) FindSalesTaxRates() entity.SalesTaxRates {
	date, _ := time.Parse("2006-01-02", "2019-10-01")
	return entity.SalesTaxRates{
		{ID: 1, Rate: decimal.NewFromFloat(0.1), AppliedDate: date, Classification: "10%"},
	}
}
