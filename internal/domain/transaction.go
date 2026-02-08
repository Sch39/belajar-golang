// internal/domain/transaction.go
package domain

type Transaction struct {
	ID         string
	TotalPrice int64

	Timestamp
	SoftDelete
}

type TransactionDetail struct {
	ID            string
	TransactionID string
	ProductID     string
	Quantity      int32
	Price         int64

	Timestamp
	SoftDelete
}
