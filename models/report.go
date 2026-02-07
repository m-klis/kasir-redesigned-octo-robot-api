package models

type TodayReport struct {
	TotalRevenue     int                `json:"total_revenue"`
	TotalTransaction int                `json:"total_transaksi"`
	BestSelling      BestSellingProduct `json:"produk_terlaris"`
}

type BestSellingProduct struct {
	Name     string `json:"nama"`
	Quantity int    `json:"qty_terjual"`
}
