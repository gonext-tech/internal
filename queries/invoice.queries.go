package queries

type InvoiceQueryParams struct {
	SearchTerm  string `query:"searchTerm"`
	Status      string `query:"status"`
	PaymentType string `query:"payment_type"`
	InvoiceType string `query:"invoice_type"`
	InvoiceDate string `query:"invoice_date"`
	CreatedBy   int    `query:"created_by"`
	PaidTo      int    `query:"paid_to"`
	Limit       int    `query:"limit"`
	Page        int    `query:"page"`
	SortBy      string `query:"sortBy"`
	OrderBy     string `query:"orderBy"`
}

func (q *InvoiceQueryParams) SetDefaults() {
	if q.SortBy == "" {
		q.SortBy = "id"
	}
	if q.OrderBy == "" {
		q.OrderBy = "desc"
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 {
		q.Limit = 20
	}
}

func (q *InvoiceQueryParams) SearchDefaults() {
	q.Status = "ACTIVE"
	q.OrderBy = "desc"
	q.SortBy = "id"
	q.Page = 1
	q.Limit = 50
}
