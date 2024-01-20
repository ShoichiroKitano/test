package entity

import (
	"github.com/shopspring/decimal"
	"sort"
	"time"
)

type SalesTaxRate struct {
	ID             uint64
	Rate           decimal.Decimal
	AppliedDate    time.Time
	Classification string
}

func (rate SalesTaxRate) SalesTax(amount int64) int64 {
	salesTaxAsDceimal := decimal.NewFromInt(amount).Mul(rate.Rate)
	salesTax, _ := salesTaxAsDceimal.BigFloat().Int64() // TODO: 端数の処理
	return salesTax
}

type SalesTaxRates []*SalesTaxRate

func (rates SalesTaxRates) FindByID(id uint64) *SalesTaxRate {
	for _, rate := range rates {
		if rate.ID == id {
			return rate
		}
	}
	return nil
}

func (rates SalesTaxRates) AppliedAt(date time.Time) *SalesTaxRate {
	tmp := make(SalesTaxRates, 0, len(rates))
	for _, rate := range rates {
		tmp = append(tmp, rate)
	}
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].AppliedDate.Before(time.Now())
	})
	for _, rate := range tmp {
		if date.Equal(rate.AppliedDate) || date.After(rate.AppliedDate) {
			return rate
		}
	}
	return nil
}
