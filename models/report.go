package models

// GET /api/report/hari-ini
// Response:
// {
//   "total_revenue": 45000,
//   "total_transaksi": 5,
//   "produk_terlaris": { "nama": "Indomie Goreng", "qty_terjual": 12 }
// }

type TodayReport struct {
	TotalRevenue     int                `json:"total_revenue"`
	TotalTransaction int                `json:"total_transaksi"`
	BestSelling      BestSellingProduct `json:"produk_terlaris"`
}

type BestSellingProduct struct {
	Name     string `json:"nama"`
	Quantity int    `json:"qty_terjual"`
}
