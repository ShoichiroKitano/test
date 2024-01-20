package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"super-shiharai-kun/entity"
	"super-shiharai-kun/repository"
	"time"
)

type InvoiceController struct {
	RDB *repository.RDBRepository
}

func (controller *InvoiceController) Index(c echo.Context) error {
	userID, err := UserID(c)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	corporate, err := controller.RDB.FindCorporateByUserID(userID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	invoices, err := controller.RDB.FindInvoiceWithPaymentDueDatesInPeriod(
		corporate.ID,
		c.QueryParam("start_date"),
		c.QueryParam("end_date"),
	)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	taxRates := controller.RDB.FindSalesTaxRates()

	invoiceJSONs := make([]map[string]any, 0, len(invoices))
	for _, invoice := range invoices {
		taxRate := taxRates.FindByID(invoice.SalesTaxRateID)
		invoiceJSONs = append(invoiceJSONs, toJSONInvoice(invoice, taxRate))
	}
	return c.JSON(http.StatusOK, map[string]any{"invoices": invoiceJSONs})
}

type InvoiceCreateRequestJson struct {
	PartnerID      uint64 `json:"partner_id"`
	IssueDate      string `json:"issue_date"` // TODO: ちゃんとした日付の型を作る
	PaymentAmount  int64  `json:"payment_amount"`
	PaymentDueDate string `json:"payment_due_date"` // TODO: ちゃんとした日付の型を作る
}

func (controller *InvoiceController) Create(c echo.Context) error {
	json := &InvoiceCreateRequestJson{}
	c.Bind(json)
	userID, err := UserID(c)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	corporate, err := controller.RDB.FindCorporateByUserID(userID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	partners, err := controller.RDB.FindPartnerByCorporateID(corporate.ID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	partner := partners.FindByPartnerID(json.PartnerID)
	if partner == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	taxRates := controller.RDB.FindSalesTaxRates()
	salesTaxRate := taxRates.AppliedAt(time.Now()) // TODO: 本来は振込時に税率を適用すべきか？（サービスの内容的には振り込みに対して手数料が発生している）
	invoice := &entity.Invoice{
		CorporateID:    corporate.ID,
		PartnerID:      partner.ID,
		IssueDate:      parseDate(json.IssueDate),
		PaymentAmount:  json.PaymentAmount,
		PaymentDueDate: parseDate(json.PaymentDueDate),
		Status:         entity.InvoiceStatusUnprocessed,
	}
	invoice.FillAmountDueAndAmountDueDetails(entity.PaymentFeeRate, salesTaxRate)
	if err = controller.RDB.CreateInvoice(invoice); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, toJSONInvoice(invoice, salesTaxRate))
}

func toJSONInvoice(invoice *entity.Invoice, salesTaxRate *entity.SalesTaxRate) map[string]any {
	var invoiceStatus string
	switch invoice.Status {
	case entity.InvoiceStatusUnprocessed:
		invoiceStatus = "未処理"
	case entity.InvoiceStatusProcessing:
		invoiceStatus = "処理中"
	case entity.InvoiceStatusPaid:
		invoiceStatus = "支払い済み"
	case entity.InvoiceStatusError:
		invoiceStatus = "エラー"
	}

	return map[string]any{
		"id":                invoice.ID,
		"corporate_id":      invoice.CorporateID,
		"partner_id":        invoice.PartnerID,
		"payment_amount":    invoice.PaymentAmount,
		"amount_due":        invoice.AmountDue,
		"fee":               invoice.Fee,
		"fee_rate":          invoice.FeeRate,
		"sales_tax_rate_id": invoice.SalesTaxRateID,
		"sales_tax_rate":    salesTaxRate.Rate,
		"status":            invoiceStatus,
		"issue_date":        toStringDate(invoice.IssueDate),
		"payment_due_date":  toStringDate(invoice.PaymentDueDate),
	}
}

func parseDate(str string) time.Time {
	date, _ := time.Parse("2006-01-02", str)
	return date
}

func toStringDate(date time.Time) string {
	year, month, day := date.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}
