package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

const (
	InvoiceStatusUnprocessed uint8 = 0
	InvoiceStatusProcessing  uint8 = 1
	InvoiceStatusCompleted   uint8 = 2
	InvoiceStatusError       uint8 = 3
)

type Invoice struct {
	ID             uint64          `db:"id"`
	CorporateID    uint64          `db:"corporate_id"`
	PartnerID      uint64          `db:"partner_id"`
	IssueDate      time.Time       `db:"issue_date"` // TODO: ちゃんとした日付の型を作る
	PaymentAmount  int64           `db:"payment_amount"`
	Fee            int64           `db:"fee"`
	FeeRate        decimal.Decimal `db:"fee_rate"`
	SalesTaxRateID uint64          `db:"sales_tax_rate_id"` // 税区分の仕様上同じ税率でも区別する必要があるのでIDを持つように仕様変更
	AmountDue      int64           `db:"amount_due"`
	PaymentDueDate time.Time       `db:"payment_due_date"`
	Status         uint8           `db:"status"`
	//BillingID int //TODO: 取引先が発行している請求番号があるはずなのでそれを使って二重振り込みを防ぐ必要あり？
}

func (invoice *Invoice) FillAmountDueAndAmountDueDetails(feeRate *FeeRate, salesTaxRate *SalesTaxRate) {
	fee := feeRate.Fee(invoice.PaymentAmount)
	salesTaxOfFee := salesTaxRate.SalesTax(fee)
	invoice.AmountDue = invoice.PaymentAmount + fee + salesTaxOfFee
	invoice.Fee = fee
	invoice.FeeRate = feeRate.Rate
	invoice.SalesTaxRateID = salesTaxRate.ID
}
