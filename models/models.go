package models

type TarifBPJSRawatInap struct {
	KodeINA   string `gorm:"column:Kode_INA_CBG"`
	Deskripsi string `gorm:"column:Dekripsi"`
	Kelas3    int    `gorm:"column:Kelas_3"`
	Kelas2    int    `gorm:"column:Kelas_2"`
	Kelas1    int    `gorm:"column:Kelas_1"`
}

type TarifBPJSRawatJalan struct {
	KodeINA     string `gorm:"column:Kode_INA_CBG"`
	Deskripsi   string `gorm:"column:Dekripsi"`
	TarifINACBG string `gorm:"column:Tarif_INA_CBG" json:"tarif_inacbg"`
}

type TarifRS struct {
	KodeRS    string `gorm:"column:Kode"`
	Deskripsi string `gorm:"column:Deskripsi"`
	Harga     int    `gorm:"column:Harga"`
	Kategori  string `gorm:"column:Kategori"`
}

func (TarifBPJSRawatJalan) TableName() string {
	return "tarif_bpjs_rawat_jalan"
}

func (TarifBPJSRawatInap) TableName() string {
	return "tarif_bpjs_rawat_inap"
}

func (TarifRS) TableName() string {
	return "tarif_rs"
}
