package services

import (
	"gorm.io/gorm"
	"main.go/models"
)

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
		}

		result = append(result, item)
	}

	return result, nil
}
