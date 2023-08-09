// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AddProductStock struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type CartItem struct {
	Product  *Product `json:"product"`
	Quantity int      `json:"quantity"`
}

type NewProduct struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type ShoppingCart struct {
	ID    string      `json:"id"`
	Items []*CartItem `json:"items"`
}