package services

import (
	"main.go/database"
	"main.go/models"

	"gorm.io/gorm"
)

// Get tarif BPJS Rawat Inap
func GetTarifBPJSRawatInap() ([]models.TarifBPJSRawatInap, error) {
	var data []models.TarifBPJSRawatInap
	if err := database.DB.Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func GetTarifBPJSRawatInapByKode(kode string) (*models.TarifBPJSRawatInap, error) {
	var data models.TarifBPJSRawatInap
	if err := database.DB.Where("Kode_INA_CBG = ?", kode).First(&data).Error; err != nil {
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
	if err := database.DB.Where("Kode_INA_CBG = ?", kode).First(&data).Error; err != nil {
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
	if err := database.DB.Where("Kode = ?", kode).First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func GetTarifRSByKategori(kategori string) ([]models.TarifRS, error) {
	var data []models.TarifRS
	if err := database.DB.Where("Kategori = ?", kategori).Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func IsNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
