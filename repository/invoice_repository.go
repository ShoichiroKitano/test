package repository

import (
	"super-shiharai-kun/entity"
)

func (repo *RDBRepository) CreateInvoice(e *entity.Invoice) error {
	_, err := repo.Driver.Exec(
		`INSERT INTO invoices (
			 corporate_id,
			 partner_id,
			 issue_date,
			 payment_amount,
			 fee,
			 fee_rate,
			 sales_tax_rate_id,
			 amount_due,
			 payment_due_date,
			 status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
		e.CorporateID,
		e.PartnerID,
		e.IssueDate,
		e.PaymentAmount,
		e.Fee,
		e.FeeRate,
		e.SalesTaxRateID,
		e.AmountDue,
		e.PaymentDueDate,
		e.Status,
	)
	// TODO: idをentityに入れる
	return err
}

func (repo *RDBRepository) FindInvoiceWithPaymentDueDatesInPeriod(startDate, endDate string) ([]*entity.Invoice, error) {
	rows, err := repo.Driver.Queryx("SELECT * FROM invoices where payment_due_date BETWEEN ? AND ?;", startDate, endDate)
	if err != nil {
		return nil, err
	}
	results := []*entity.Invoice{}
	for rows.Next() {
		e := &entity.Invoice{}
		if err := rows.StructScan(e); err != nil {
			return nil, err
		}
		results = append(results, e)
	}
	return results, nil
}