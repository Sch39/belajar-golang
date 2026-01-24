// internal\category\request.go
package category

type upsertRequest struct {
	Name        string `json:"name" validate:"required,min=3" example:"Coffe"`
	Description string `json:"description"`
}
