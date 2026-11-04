package domain

type ProductDetail struct {
	Product
	Description string       `json:"description"`
	About       string       `json:"about"`
	Offers      []StoreOffer `json:"offers"`
	History     []PricePoint `json:"history"`
}

type StoreOffer struct {
	Store    Store   `json:"store"`
	Price    float64 `json:"price"`
	OldPrice float64 `json:"old_price"`
	InStock  bool    `json:"in_stock"`
}

type PricePoint struct {
	Day   string  `json:"day"`
	Price float64 `json:"price"`
}
