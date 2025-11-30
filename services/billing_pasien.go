package services

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"main.go/database"
	"main.go/models"
)




// Get ID_tarif_RS dari Tindakan_RS
func GetTarifRSByTindakan(tindakans []string) ([]models.TarifRS, error) {
	var tarifList []models.TarifRS

	if err := database.DB.
		Where("Tindakan_RS IN ?", tindakans).
		Find(&tarifList).Error; err != nil {
		return nil, err
	}

	return tarifList, nil
}



// GetPasienByID mencari pasien berdasarkan ID
func GetPasienByID(id int) (*models.Pasien, error) {
	var pasien models.Pasien

	if err := database.DB.Where("ID_Pasien = ?", id).First(&pasien).Error; err != nil {
		return nil, err
	}

	return &pasien, nil
}

// GetPasienByNama mencari pasien berdasarkan nama
func GetPasienByNama(nama string) (*models.Pasien, error) {
	var pasien models.Pasien

	if err := database.DB.Where("Nama_Pasien = ?", nama).First(&pasien).Error; err != nil {
		return nil, err
	}

	return &pasien, nil
}
 
// GetDokterByNama mencari dokter berdasarkan nama
func GetDokterByNama(nama string) (*models.Dokter, error) {
	var dokter models.Dokter

	if err := database.DB.Where("Nama_Dokter = ?", nama).First(&dokter).Error; err != nil {
		return nil, err
	}

	return &dokter, nil
} 

func DataFromFE(input models.BillingRequest) (
	*models.BillingPasien,
	*models.Pasien,
	[]models.Billing_Tindakan,
	[]models.Billing_ICD9,
	[]models.Billing_ICD10,
	error,
) {

	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, nil, nil, nil, nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ===========================
	// 1. CARI ATAU BUAT PASIEN
	// ===========================
	var pasien models.Pasien
	result := tx.Where("Nama_Pasien = ?", input.Nama_Pasien).First(&pasien)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {

			pasien = models.Pasien{
				Nama_Pasien:   input.Nama_Pasien,
				Jenis_Kelamin: input.Jenis_Kelamin,
				Usia:          input.Usia,
				Ruangan:       input.Ruangan,
				Kelas:         input.Kelas,
			}

			if err := tx.Create(&pasien).Error; err != nil {
				tx.Rollback()
				return nil, nil, nil, nil, nil,
					fmt.Errorf("gagal membuat pasien baru: %s", err.Error())
			}
		} else {
			tx.Rollback()
			return nil, nil, nil, nil, nil,
				fmt.Errorf("gagal mencari pasien: %s", result.Error.Error())
		}
	}

	if pasien.ID_Pasien == 0 {
		tx.Rollback()
		return nil, nil, nil, nil, nil, fmt.Errorf("ID_Pasien tidak valid")
	}

	// ===========================
	// 2. CARI DOKTER
	// ===========================
	var dokter models.Dokter
	if err := tx.Where("Nama_Dokter = ?", input.Nama_Dokter).First(&dokter).Error; err != nil {
		tx.Rollback()
		return nil, nil, nil, nil, nil,
			fmt.Errorf("dokter '%s' tidak ditemukan", input.Nama_Dokter)
	}

	// ===========================
	// 3. BUAT BILLING
	// ===========================
	now := time.Now()
	billing := models.BillingPasien{
		ID_Pasien:        pasien.ID_Pasien,
		Cara_Bayar:       input.Cara_Bayar,
		Tanggal_masuk:    &now,
		ID_Dokter:        dokter.ID_Dokter,
		Total_Tarif_RS:   input.Total_Tarif_RS,
		Total_Tarif_BPJS: 0,
		Billing_sign:     "created",
	}

	if err := tx.Create(&billing).Error; err != nil {
		tx.Rollback()
		return nil, nil, nil, nil, nil,
			fmt.Errorf("gagal membuat billing: %s", err.Error())
	}

	// ===========================================================
	// 4. PROSES TINDAKAN_RΣ, ICD9, ICD10
	// ===========================================================

	var billingTindakanList []models.Billing_Tindakan
	var billingICD9List []models.Billing_ICD9
	var billingICD10List []models.Billing_ICD10

	// ----------- TINDAKAN RS -----------
	for _, tindakan := range input.Tindakan_RS {
		var tarif models.TarifRS

		if err := tx.Where("Tindakan_RS = ?", tindakan).First(&tarif).Error; err != nil {
			tx.Rollback()
			return nil, nil, nil, nil, nil,
				fmt.Errorf("tindakan '%s' tidak ditemukan", tindakan)
		}

		billTindakan := models.Billing_Tindakan{
			ID_Billing:  billing.ID_Billing,
			ID_Tarif_RS: tarif.KodeRS,
		}

		if err := tx.Create(&billTindakan).Error; err != nil {
			tx.Rollback()
			return nil, nil, nil, nil, nil,
				fmt.Errorf("gagal insert billing tindakan: %s", err.Error())
		}

		billingTindakanList = append(billingTindakanList, billTindakan)
	}

	// ----------- ICD9 -----------
	for _, icd := range input.ICD9 {
		var icd9 models.ICD9

		if err := tx.Where("Prosedur = ?", icd).First(&icd9).Error; err != nil {
			tx.Rollback()
			return nil, nil, nil, nil, nil,
				fmt.Errorf("ICD9 '%s' tidak ditemukan", icd)
		}

		billICD9 := models.Billing_ICD9{
			ID_Billing: billing.ID_Billing,
			ID_ICD9:    icd9.Kode_ICD9,
		}

		if err := tx.Create(&billICD9).Error; err != nil {
			tx.Rollback()
			return nil, nil, nil, nil, nil,
				fmt.Errorf("gagal insert billing ICD9: %s", err.Error())
		}

		billingICD9List = append(billingICD9List, billICD9)
	}

	// ----------- ICD10 -----------
	for _, icd := range input.ICD10 {
		var icd10 models.ICD10

		if err := tx.Where("Diagnosa = ?", icd).First(&icd10).Error; err != nil {
			tx.Rollback()
			return nil, nil, nil, nil, nil,
				fmt.Errorf("ICD10 '%s' tidak ditemukan", icd)
		}

		billICD10 := models.Billing_ICD10{
			ID_Billing: billing.ID_Billing,
			ID_ICD10:   icd10.Kode_ICD10,
		}

		if err := tx.Create(&billICD10).Error; err != nil {
			tx.Rollback()
			return nil, nil, nil, nil, nil,
				fmt.Errorf("gagal insert billing ICD10: %s", err.Error())
		}

		billingICD10List = append(billingICD10List, billICD10)
	}

	// ===========================
	// 5. COMMIT
	// ===========================
	if err := tx.Commit().Error; err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return &billing, &pasien, billingTindakanList, billingICD9List, billingICD10List, nil
}

