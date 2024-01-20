package entity

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvoice_FillAmountDueAndAmountDueDetails(t *testing.T) {
	invoice := &Invoice{PaymentAmount: 1000}
	invoice.FillAmountDueAndAmountDueDetails(
		&FeeRate{Rate: decimal.NewFromFloat(0.04)},
		&SalesTaxRate{ID: 1, Rate: decimal.NewFromFloat(0.1)},
	)
	assert.Equal(t, int64(1044), invoice.AmountDue)
	assert.Equal(t, int64(40), invoice.Fee)
	assert.Equal(t, decimal.NewFromFloat(0.04), invoice.FeeRate)
	assert.Equal(t, uint64(1), invoice.SalesTaxRateID)
}
