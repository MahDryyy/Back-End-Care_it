package models

type TarifBPJSRawatInap struct {
	KodeINA   string  `gorm:"column:ID_INACBG_RI"`
	Deskripsi string  `gorm:"column:Tindakan_RI"`
	Kelas1    float64 `gorm:"column:Tarif_Kelas_1"`
	Kelas2    float64 `gorm:"column:Tarif_Kelas_2"`
	Kelas3    float64 `gorm:"column:Tarif_Kelas_3"`
}

type TarifBPJSRawatJalan struct {
	KodeINA     string  `gorm:"column:ID_INACBG_RJ"`
	Deskripsi   string  `gorm:"column:Tindakan_RJ"`
	TarifINACBG float64 `gorm:"column:Tarif_RJ" json:"tarif_inacbg"`
}

type TarifRS struct {
	KodeRS    string `gorm:"column:ID_Tarif_RS"`
	Deskripsi string `gorm:"column:Tindakan_RS"`
	Harga     int    `gorm:"column:Tarif_RS"`
	Kategori  string `gorm:"column:Kategori_RS"`
}

func (TarifBPJSRawatJalan) TableName() string {
	return "ina_cbg_rawatjalan"
}

func (TarifBPJSRawatInap) TableName() string {
	return "ina_cbg_rawatinap"
}

func (TarifRS) TableName() string {
	return "tarif_rs"
}

// ICD9

type ICD9 struct {
	Kode_ICD9 string `gorm:"column:ID_ICD9"`
	Prosedur  string `gorm:"column:Prosedur"`
	Versi     string `gorm:"column:Versi_ICD9"`
}

func (ICD9) TableName() string {
	return "icd9"
}

// ICD10
type ICD10 struct {
	Kode_ICD10 string `gorm:"column:ID_ICD10"`
	Diagnosa   string `gorm:"column:Diagnosa"`
	Versi      string `gorm:"column:Versi_ICD10"`
}

func (ICD10) TableName() string {
	return "icd10"
}

// ruangan
type Ruangan struct {
	ID_Ruangan       string `gorm:"column:ID_Ruangan"`
	Jenis_Ruangan    string `gorm:"column:Jenis_Ruangan"`
	Nama_Ruangan     string `gorm:"column:Nama_Ruangan"`
	Keterangan       string `gorm:"column:keterangan"`
	Kategori_ruangan string `gorm:"column:kategori_ruangan"`
}

func (Ruangan) TableName() string {
	return "ruangan"
}

// dokter
type Dokter struct {
	ID_Dokter     string `gorm:"column:ID_Dokter"`
	Nama_Dokter   string `gorm:"column:Nama_Dokter"`
	Password      string `gorm:"column:Password"`
	Status        string `gorm:"column:Status"`
	KSM           string `gorm:"column:KSM"`
	Email_UB      string `gorm:"column:Email_UB"`
	Email_Pribadi string `gorm:"column:Email_Pribadi"`
}

func (Dokter) TableName() string {
	return "dokter"
}
