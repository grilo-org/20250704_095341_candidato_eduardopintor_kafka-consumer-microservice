package dto

type Order struct {
	OrderID     string `json:"orderDeVenda"`
	OrderStatus string `json:"estapaAtual"`
}
