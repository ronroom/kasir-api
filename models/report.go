package models

type BestSellingProd struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

type SalesSummary struct {
	TotalRevenue   int             `json:"total_revenue"`
	TotalTransaksi int             `json:"total_transaksi"`
	ProdukTerlaris BestSellingProd `json:"produk_terlaris"`
}
