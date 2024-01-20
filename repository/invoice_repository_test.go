package repository

import (
	"errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"super-shiharai-kun/driver"
	"super-shiharai-kun/entity"
	"testing"
	"time"
)

func TestInvoiceRepository_FindInvoiceByIssueDates(t *testing.T) {
	RunWithTx(t, "永続化されたInvoiceが取得できること", func(t *testing.T, tx RDBRepository) {
		require.NoError(t, tx.CreateInvoice(&entity.Invoice{
			FeeRate:        decimal.NewFromFloat(0.10),
			Fee:            10,
			PaymentAmount:  100,
			AmountDue:      1000,
			CorporateID:    2,
			PartnerID:      3,
			SalesTaxRateID: 4,
			PaymentDueDate: parseDate("2024-01-01"),
			IssueDate:      parseDate("2024-01-02"),
			Status:         entity.InvoiceStatusPaid,
		}))

		results, err := tx.FindInvoiceWithPaymentDueDatesInPeriod(2, "2024-01-01", "2024-01-01")
		require.NoError(t, err)
		require.Equal(t, 1, len(results))
		assert.NotEqual(t, uint64(0), results[0].ID)
		assert.Equal(t, decimal.NewFromFloat(0.10).BigFloat(), results[0].FeeRate.BigFloat())
		assert.Equal(t, int64(10), results[0].Fee)
		assert.Equal(t, int64(100), results[0].PaymentAmount)
		assert.Equal(t, int64(1000), results[0].AmountDue)
		assert.Equal(t, uint64(2), results[0].CorporateID)
		assert.Equal(t, uint64(3), results[0].PartnerID)
		assert.Equal(t, uint64(4), results[0].SalesTaxRateID)
		assert.Equal(t, parseDate("2024-01-01"), results[0].PaymentDueDate)
		assert.Equal(t, parseDate("2024-01-02"), results[0].IssueDate)
		assert.Equal(t, entity.InvoiceStatusPaid, results[0].Status)
	})

	RunWithTx(t, "指定した期間のInvoiceが取得できること", func(t *testing.T, tx RDBRepository) {
		require.NoError(t, tx.CreateInvoice(createInvoice(map[string]any{"PaymentDueDate": parseDate("2023-12-31"), "CorporateID": 1})))
		require.NoError(t, tx.CreateInvoice(createInvoice(map[string]any{"PaymentDueDate": parseDate("2024-01-01"), "CorporateID": 1})))
		require.NoError(t, tx.CreateInvoice(createInvoice(map[string]any{"PaymentDueDate": parseDate("2024-01-02"), "CorporateID": 1})))
		require.NoError(t, tx.CreateInvoice(createInvoice(map[string]any{"PaymentDueDate": parseDate("2024-01-03"), "CorporateID": 1})))
		require.NoError(t, tx.CreateInvoice(createInvoice(map[string]any{"PaymentDueDate": parseDate("2024-01-04"), "CorporateID": 1})))

		results, err := tx.FindInvoiceWithPaymentDueDatesInPeriod(1, "2024-01-01", "2024-01-03")
		require.NoError(t, err)
		require.Equal(t, 3, len(results))
		assert.Equal(t, parseDate("2024-01-01"), results[0].PaymentDueDate)
		assert.Equal(t, parseDate("2024-01-02"), results[1].PaymentDueDate)
		assert.Equal(t, parseDate("2024-01-03"), results[2].PaymentDueDate)
	})

	RunWithTx(t, "他の企業のinvoiceは取得できないこと", func(t *testing.T, tx RDBRepository) {
		invoiceCreatedCorporateID := uint64(1)
		require.NoError(t, tx.CreateInvoice(
			createInvoice(map[string]any{"PaymentDueDate": parseDate("2024-01-01"), "CorporateID": invoiceCreatedCorporateID})),
		)

		otherCorporateID := uint64(2)
		results, err := tx.FindInvoiceWithPaymentDueDatesInPeriod(otherCorporateID, "2024-01-01", "2024-01-01")
		require.NoError(t, err)
		require.Equal(t, 0, len(results))
	})
}

// TODO: テストヘルパーにする
func createInvoice(override map[string]any) *entity.Invoice {
	invoice := &entity.Invoice{
		PaymentDueDate: parseDate("1970-01-01"),
		IssueDate:      parseDate("1970-01-01"),
	}
	if value, ok := override["PaymentDueDate"]; ok {
		invoice.PaymentDueDate = value.(time.Time)
	}
	if value, ok := override["CorporateID"]; ok {
		switch id := value.(type) {
		case int:
			invoice.CorporateID = uint64(id)
		case uint64:
			invoice.CorporateID = id
		}
	}
	return invoice
}

// TODO: テストヘルパーにする
func parseDate(str string) time.Time {
	date, _ := time.Parse("2006-01-02", str)
	return date
}

func RunWithTx(t *testing.T, name string, f func(t *testing.T, tx RDBRepository)) {
	driver, err := driver.NewRDBDriver() //TODO: テスト間でコネクションを共有する
	require.NoError(t, err)
	repo := &RDBRepository{Driver: driver}

	t.Run(name, func(t *testing.T) {
		err := repo.Tx(func(tx RDBRepository) error {
			f(t, tx)
			return errors.New("(rollback)")
		})
		if err.Error() != "(rollback)" {
			require.NoError(t, err)
		}
	})
}
