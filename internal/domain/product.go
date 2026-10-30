package domain

type Product struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Variant  string  `json:"variant,omitempty"`
	Size     string  `json:"size,omitempty"`
	Barcode  string  `json:"barcode,omitempty"`
	ImageURL string  `json:"image_url,omitempty"`
	Price    float64 `json:"price,omitempty"`
	Store    string  `json:"store,omitempty"`
}

type PriceUpdate struct {
	ProductName string  `json:"product_name"`
	StoreName   string  `json:"store_name"`
	Price       float64 `json:"price"`
	Notes       string  `json:"notes,omitempty"`
	PhotoURL    string  `json:"photo_url,omitempty"`
}
