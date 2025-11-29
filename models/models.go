package models

import "time"

// Tarif Models

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
	ID_Dokter     int    `gorm:"column:ID_Dokter;primaryKey"`
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

// PASIEN
type Pasien struct {
	ID_Pasien     int    `gorm:"column:ID_Pasien;primaryKey;autoIncrement"`
	Nama_Pasien   string `gorm:"column:Nama_Pasien"`
	Jenis_Kelamin string `gorm:"column:Jenis_Kelamin"`
	Usia          int    `gorm:"column:Usia"`
	Ruangan       string `gorm:"column:Ruangan"`
	Kelas         string `gorm:"column:Kelas"`
}

type Kelas string

const (
	Kelas_1 Kelas = "1"
	Kelas_2 Kelas = "2"
	Kelas_3 Kelas = "3"
)

type Jenis_kelamin string

const (
	Jenis_Kelamin_Laki_laki Jenis_kelamin = "Laki-laki"
	Jenis_Kelamin_Perempuan Jenis_kelamin = "Perempuan"
)

func (Pasien) TableName() string {
	return "pasien"
}

//billing pasien

type BillingPasien struct {
	ID_Billing       int        `gorm:"column:ID_Billing;primaryKey;autoIncrement"`
	ID_Pasien        int        `gorm:"column:ID_Pasien"`
	Cara_Bayar       string     `gorm:"column:Cara_Bayar"`
	Tanggal_masuk    *time.Time `gorm:"column:Tanggal_Masuk"`
	Tanggal_keluar   *time.Time `gorm:"column:Tanggal_Keluar"`
	ID_Dokter        int        `gorm:"column:ID_Dokter"`
	Total_Tarif_RS   int        `gorm:"column:Total_Tarif_RS"`
	Total_Tarif_BPJS float64    `gorm:"column:Total_klaim"`
	Billing_sign     string     `gorm:"column:Billing_sign"`
}

type Billing_sign string

const (
	Billing_Sign_hijau  Billing_sign = "hijau"
	Billing_Sign_kuning Billing_sign = "kuning"
	Billing_Sign_merah  Billing_sign = "merah"
)

type Cara_bayar string

const (
	Cara_Bayar_BPJS Cara_bayar = "BPJS"
	Cara_Bayar_UMUM Cara_bayar = "UMUM"
)

func (BillingPasien) TableName() string {
	return "billing_pasien"
}

// BillingRequest untuk menerima data dari API
type BillingRequest struct {
	Nama_Dokter    string `json:"nama_dokter" binding:"required"`
	Nama_Pasien    string `json:"nama_pasien" binding:"required"`
	Jenis_Kelamin  string `json:"jenis_kelamin" binding:"required"`
	Usia           int    `json:"usia" binding:"required"`
	Ruangan        string `json:"ruangan" binding:"required"`
	Kelas          string `json:"kelas" binding:"required"`
	Tindakan_RS    string `json:"tindakan_rs" binding:"required"`
	ICD9           string `json:"icd9" binding:"required"`
	ICD10          string `json:"icd10" binding:"required"`
	Cara_Bayar     string `json:"cara_bayar" binding:"required"`
	Total_Tarif_RS int    `json:"total_tarif_rs"`
}

// admin ruangan

type admin_ruangan struct {
	ID_Admin   int    `gorm:"column:ID_Admin"`
	Nama_Admin string `gorm:"column:Nama_Admin"`
	Password   string `gorm:"column:Password"`
	ID_Ruangan string `gorm:"column:ID_Ruangan"`
}

func (admin_ruangan) TableName() string {
	return "admin_ruangan"
}

// inputan dari FE
