package entity

import (
	"github.com/shopspring/decimal"
)

var PaymentFeeRate *FeeRate = &FeeRate{Rate: decimal.NewFromFloat(0.04)} //TODO: 仮実装

// TODO: 適用日等を設定してDBで管理する
type FeeRate struct {
	Rate decimal.Decimal
}

func (feeRate *FeeRate) Fee(amount int64) int64 {
	feeAsDecimal := feeRate.Rate.Mul(decimal.NewFromInt(amount))
	fee, _ := feeAsDecimal.BigFloat().Int64() // TODO: 端数の処理
	return fee
}
