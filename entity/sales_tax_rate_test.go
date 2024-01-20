package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSalesTaxRates_AppliedAt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		subject                string
		rates                  SalesTaxRates
		inputDate              time.Time
		expectedSalesTaxRateID uint64
	}{
		{
			subject: "指定した日付がより未来の適用日の消費税率を取得する",
			rates: SalesTaxRates{
				{ID: 1, AppliedDate: parseDate("2023-01-01")},
				{ID: 2, AppliedDate: parseDate("2023-02-01")},
			},
			inputDate:              parseDate("2023-02-02"),
			expectedSalesTaxRateID: 2,
		},
		{
			subject: "指定したより未来の適用日の消費税率が複数ある場合は、最も古い日付の消費税率を取得する",
			rates: SalesTaxRates{
				{ID: 1, AppliedDate: parseDate("2023-01-01")},
				{ID: 3, AppliedDate: parseDate("2023-03-01")},
				{ID: 2, AppliedDate: parseDate("2023-02-01")},
			},
			inputDate:              parseDate("2023-02-02"),
			expectedSalesTaxRateID: 2,
		},
		{
			subject: "日付が一致する場合",
			rates: SalesTaxRates{
				{ID: 2, AppliedDate: parseDate("2023-01-02")},
				{ID: 1, AppliedDate: parseDate("2023-01-01")},
			},
			inputDate:              parseDate("2023-01-01"),
			expectedSalesTaxRateID: 1,
		},
	}
	for _, tt := range tests {
		result := tt.rates.AppliedAt(tt.inputDate)
		assert.Equal(t, tt.expectedSalesTaxRateID, result.ID)
	}
}

func parseDate(str string) time.Time {
	date, _ := time.Parse("2006-01-02", str)
	return date
}
