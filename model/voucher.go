package model

type Voucher struct {
	VoucherID            uint
	Name                 string `json:"name"`
	KodeVoucher          string `json:"kode_voucher"`
	JumlahPotonganPersen int    `json:"jumlah_potongan_persen"`
}
