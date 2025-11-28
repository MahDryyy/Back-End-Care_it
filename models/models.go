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

// ICD9

type ICD9 struct {
	Kode_ICD9 string `gorm:"column:ID_ICD9"`
	Prosedur string `gorm:"column:Prosedur"`
	Versi string `gorm:"column:Versi_ICD9"`
}

func (ICD9) TableName() string {
	return "icd9"
}

// ICD10
type ICD10 struct {
	Kode_ICD10 string `gorm:"column:ID_ICD10"`
	Diagnosa string `gorm:"column:Diagnosa"`
	Versi string `gorm:"column:Versi_ICD10"`
}

func (ICD10) TableName() string {
	return "icd10"
}

// ruangan
type Ruangan struct {
	ID_Ruangan string `gorm:"column:ID_Ruangan"`
	Jenis_Ruangan string `gorm:"column:Jenis_Ruangan"`
	Nama_Ruangan string `gorm:"column:Nama_Ruangan"`

}

func (Ruangan) TableName() string {
	return "ruangan"
}

// dokter
type Dokter struct {
	ID_Dokter string `gorm:"column:ID_Dokter"`
	Nama_Dokter string `gorm:"column:Nama_Dokter"`
	Password string `gorm:"column:Password"`
	Status string `gorm:"column:Status"`
	KSM string `gorm:"column:KSM"`
	Email_UB string `gorm:"column:Email_UB"`
	Email_Pribadi string `gorm:"column:Email_Pribadi"`
}

func (Dokter) TableName() string {
	return "dokter"
}