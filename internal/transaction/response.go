// internal/transaction/response.go
package transaction

type checkoutResponse struct {
	ID         string `json:"id" example:"9825b44a-101f-4c6e-8d8a-675204481359"`
	TotalPrice int64  `json:"total_price" example:"1000000"`
}

// type reportResponse struct {
// 	TotalRevenue       int64 `json:"total_revenue" example:"1000000"`
// 	TotalTransaction   int64 `json:"total_transaction" example:"100"`
// 	BestSellingProduct struct {
// 		ID      string `json:"id" example:"9825b44a-101f-4c6e-8d8a-675204481359"`
// 		Name    string `json:"name" example:"Product 1"`
// 		QtySold int64  `json:"qty_sold" example:"100"`
// 	}
// }
