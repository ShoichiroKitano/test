package repository

import (
	"super-shiharai-kun/entity"
)

func (repo *RDBRepository) CreateInvoice(e *entity.Invoice) error {
	result, err := repo.Driver.Exec(
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
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	e.ID = uint64(id)
	return err
}

func (repo *RDBRepository) FindInvoiceWithPaymentDueDatesInPeriod(corporateID uint64, startDate, endDate string) ([]*entity.Invoice, error) {
	rows, err := repo.Driver.Queryx(
		"SELECT * FROM invoices where corporate_id = ? and payment_due_date BETWEEN ? AND ?;",
		corporateID,
		startDate,
		endDate,
	)
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
