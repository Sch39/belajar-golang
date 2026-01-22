// internal\product\response.go
package product

type productResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}