package services

import (
	"errors"
	"fmt"

	"backendcareit/models"

	"gorm.io/gorm"
)

func Post_INACBG_Admin(db *gorm.DB, input models.Post_INACBG_Admin) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Ensure rollback on panic / unexpected error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Validate input
	if input.Tipe_inacbg != "RI" && input.Tipe_inacbg != "RJ" {
		tx.Rollback()
		return errors.New("invalid tipe_inacbg: must be 'RI' or 'RJ'")
	}
	if len(input.Kode_INACBG) == 0 {
		tx.Rollback()
		return errors.New("Kode_INACBG tidak boleh kosong")
	}

	// 1. Update total klaim dan billing_sign
	res := tx.Model(&models.BillingPasien{}).
		Where("ID_Billing = ?", input.ID_Billing).
		Updates(map[string]interface{}{
			"Total_klaim":  input.Total_klaim,
			"Billing_sign": input.Billing_sign,
		})

	if res.Error != nil {
		tx.Rollback()
		return fmt.Errorf("gagal update billing: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("billing dengan ID_Billing=%d tidak ditemukan", input.ID_Billing)
	}

	// 2. Bulk insert kode INACBG berdasarkan tipe_inacbg
	switch input.Tipe_inacbg {
	case "RI":
		records := make([]models.Billing_INACBG_RI, 0, len(input.Kode_INACBG))
		for _, kode := range input.Kode_INACBG {
			records = append(records, models.Billing_INACBG_RI{
				ID_Billing:  input.ID_Billing,
				Kode_INACBG: kode,
			})
		}
		if err := tx.Create(&records).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal insert INACBG RI: %w", err)
		}

	case "RJ":
		records := make([]models.Billing_INACBG_RJ, 0, len(input.Kode_INACBG))
		for _, kode := range input.Kode_INACBG {
			records = append(records, models.Billing_INACBG_RJ{
				ID_Billing:  input.ID_Billing,
				Kode_INACBG: kode,
			})
		}
		if err := tx.Create(&records).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal insert INACBG RJ: %w", err)
		}
	}

	return tx.Commit().Error
}

func GetAllBilling(db *gorm.DB) ([]models.Request_Admin_Inacbg, error) {
	var billings []models.BillingPasien

	// Ambil semua billing
	if err := db.Find(&billings).Error; err != nil {
		return nil, err
	}

	// Kumpulkan semua ID_Billing dan ID_Pasien
	var billingIDs []int
	var pasienIDs []int

	for _, b := range billings {
		billingIDs = append(billingIDs, b.ID_Billing)
		pasienIDs = append(pasienIDs, b.ID_Pasien)
	}

	// Ambil semua pasien hanya untuk pasien yang ada di billing
	pasienMap := make(map[int]models.Pasien)
	var pasienList []models.Pasien

	if err := db.Where("ID_Pasien IN ?", pasienIDs).Find(&pasienList).Error; err != nil {
		return nil, err
	}

	for _, p := range pasienList {
		pasienMap[p.ID_Pasien] = p
	}

	// Ambil tindakan hanya untuk billing terkait
	tindakanMap := make(map[int][]string)
	var tindakanRows []struct {
		ID_Billing int
		Kode       string
	}

	if err := db.Table("billing_tindakan").
		Where("ID_Billing IN ?", billingIDs).
		Select("ID_Billing, ID_Tarif_RS as Kode").
		Scan(&tindakanRows).Error; err != nil {
		return nil, err
	}

	for _, t := range tindakanRows {
		tindakanMap[t.ID_Billing] = append(tindakanMap[t.ID_Billing], t.Kode)
	}

	// Ambil ICD9
	icd9Map := make(map[int][]string)
	var icd9Rows []struct {
		ID_Billing int
		Kode       string
	}

	if err := db.Table("billing_icd9").
		Where("ID_Billing IN ?", billingIDs).
		Select("ID_Billing, ID_ICD9 as Kode").
		Scan(&icd9Rows).Error; err != nil {
		return nil, err
	}

	for _, row := range icd9Rows {
		icd9Map[row.ID_Billing] = append(icd9Map[row.ID_Billing], row.Kode)
	}

	// Ambil ICD10
	icd10Map := make(map[int][]string)
	var icd10Rows []struct {
		ID_Billing int
		Kode       string
	}

	if err := db.Table("billing_icd10").
		Where("ID_Billing IN ?", billingIDs).
		Select("ID_Billing, ID_ICD10 as Kode").
		Scan(&icd10Rows).Error; err != nil {
		return nil, err
	}

	for _, row := range icd10Rows {
		icd10Map[row.ID_Billing] = append(icd10Map[row.ID_Billing], row.Kode)
	}

	// Compile final response
	var result []models.Request_Admin_Inacbg

	for _, b := range billings {
		pasien := pasienMap[b.ID_Pasien]

		item := models.Request_Admin_Inacbg{
			ID_Billing:     b.ID_Billing,
			Nama_pasien:    pasien.Nama_Pasien,
			ID_Pasien:      pasien.ID_Pasien,
			Kelas:          pasien.Kelas,
			Ruangan:        pasien.Ruangan,
			Total_Tarif_RS: b.Total_Tarif_RS,
			Tindakan_RS:    tindakanMap[b.ID_Billing],
			ICD9:           icd9Map[b.ID_Billing],
			ICD10:          icd10Map[b.ID_Billing],
			Billing_sign:   b.Billing_sign,
		}

		result = append(result, item)
	}

	return result, nil
}
