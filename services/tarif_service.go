package services

import (
	"backendcareit/database"
	"backendcareit/models"

	"gorm.io/gorm"
)

// Get tarif BPJS Rawat Inap
func GetTarifBPJSRawatInap() ([]models.TarifBPJSRawatInap, error) {
	var data []models.TarifBPJSRawatInap
	if err := database.DB.Model(&models.TarifBPJSRawatInap{}).Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func GetTarifBPJSRawatInapByKode(kode string) (*models.TarifBPJSRawatInap, error) {
	var data models.TarifBPJSRawatInap
	if err := database.DB.Model(&models.TarifBPJSRawatInap{}).Where("ID_INACBG_RI = ?", kode).First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

// Get tarif rawat jalan
func GetTarifBPJSRawatJalan() ([]models.TarifBPJSRawatJalan, error) {
	var data []models.TarifBPJSRawatJalan
	if err := database.DB.Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func GetTarifBPJSRawatJalanByKode(kode string) (*models.TarifBPJSRawatJalan, error) {
	var data models.TarifBPJSRawatJalan
	if err := database.DB.Where("ID_INACBG_RJ = ?", kode).First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

// Get tarif RS
func GetTarifRS() ([]models.TarifRS, error) {
	var data []models.TarifRS
	if err := database.DB.Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func GetTarifRSByKode(kode string) (*models.TarifRS, error) {
	var data models.TarifRS
	if err := database.DB.Where("ID_Tarif_RS = ?", kode).First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func GetTarifRSByKategori(kategori string) ([]models.TarifRS, error) {
	var data []models.TarifRS
	if err := database.DB.Where("Kategori_RS = ?", kategori).Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func IsNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}

// ICD9
func GetICD9() ([]models.ICD9, error) {
	var data []models.ICD9
	if err := database.DB.Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

// ICD10
func GetICD10() ([]models.ICD10, error) {
	var data []models.ICD10
	if err := database.DB.Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

// ruangan
func GetRuangan() ([]models.Ruangan, error) {
	var data []models.Ruangan
	if err := database.DB.Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

// GetRuanganWithPasien - Get ruangan yang memiliki minimal 1 pasien
func GetRuanganWithPasien(db *gorm.DB) ([]models.Ruangan, error) {
	var data []models.Ruangan
	// JOIN dengan pasien table dan filter yang punya pasien
	if err := db.Distinct("ruangan.*").
		Joins("INNER JOIN pasien ON ruangan.ID_Ruangan = pasien.ID_Ruangan").
		Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

// dokter
func GetDokter() ([]models.Dokter, error) {
	var data []models.Dokter
	if err := database.DB.Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}
