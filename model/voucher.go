package model

import "gorm.io/gorm"

type Voucher struct {
	gorm.Model

	Name                 string `json:"name"`
	KodeVoucher          string `json:"kode_voucher"`
	JumlahPotonganPersen int    `json:"jumlah_potongan_persen"`
}
