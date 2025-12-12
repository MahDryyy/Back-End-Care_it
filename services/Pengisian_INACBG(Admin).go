package services

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"backendcareit/models"

	"gorm.io/gorm"
)

func Post_INACBG_Admin(db *gorm.DB, input models.Post_INACBG_Admin) error {
	// Debug log
	log.Printf("[INACBG] Input received: ID_Billing=%d, Tipe=%s, Kode_count=%d, Total_klaim=%.2f, BillingSign=%s\n",
		input.ID_Billing, input.Tipe_inacbg, len(input.Kode_INACBG), input.Total_klaim, input.Billing_sign)

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("[INACBG] Error starting transaction: %v\n", tx.Error)
		return tx.Error
	}

	// Ensure rollback on panic / unexpected error
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[INACBG] Panic recovered: %v\n", r)
			tx.Rollback()
		}
	}()

	// Validate input
	if input.Tipe_inacbg != "RI" && input.Tipe_inacbg != "RJ" {
		tx.Rollback()
		err := errors.New("invalid tipe_inacbg: must be 'RI' or 'RJ'")
		log.Printf("[INACBG] Validation error: %v\n", err)
		return err
	}
	if len(input.Kode_INACBG) == 0 {
		tx.Rollback()
		err := errors.New("Kode_INACBG tidak boleh kosong")
		log.Printf("[INACBG] Validation error: %v\n", err)
		return err
	}

	// 1. Ambil billing dulu untuk dapatkan total klaim lama
	var existingBilling models.BillingPasien
	if err := tx.Where("ID_Billing = ?", input.ID_Billing).First(&existingBilling).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("billing dengan ID_Billing=%d tidak ditemukan", input.ID_Billing)
			log.Printf("[INACBG] %v\n", err)
			return err
		}
		log.Printf("[INACBG] Error fetching billing: %v\n", err)
		return fmt.Errorf("gagal mengambil billing: %w", err)
	}

	log.Printf("[INACBG] Found billing: ID=%d, Current_Total_Klaim=%.2f\n", existingBilling.ID_Billing, existingBilling.Total_Tarif_BPJS)

	// Hitung total klaim baru = lama + tambahan
	newTotalKlaim := existingBilling.Total_Tarif_BPJS + input.Total_klaim
	log.Printf("[INACBG] New total klaim: %.2f + %.2f = %.2f\n", existingBilling.Total_Tarif_BPJS, input.Total_klaim, newTotalKlaim)

	// Parse Tanggal_Keluar jika diisi oleh admin
	var keluarPtr *time.Time
	if input.Tanggal_keluar != "" && input.Tanggal_keluar != "null" {
		s := input.Tanggal_keluar
		var parsed time.Time
		var err error
		layouts := []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"}
		for _, layout := range layouts {
			parsed, err = time.Parse(layout, s)
			if err == nil {
				t := parsed
				keluarPtr = &t
				log.Printf("[INACBG] Parsed tanggal_keluar: %v\n", t)
				break
			}
		}
		if keluarPtr == nil {
			tx.Rollback()
			err := fmt.Errorf("invalid tanggal_keluar format: %s", input.Tanggal_keluar)
			log.Printf("[INACBG] %v\n", err)
			return err
		}
	}

	// 2. Update total klaim kumulatif, billing_sign, dan tanggal keluar (jika diisi)
	updateData := map[string]interface{}{
		"Total_Klaim":  newTotalKlaim,
		"Billing_Sign": input.Billing_sign,
	}
	if keluarPtr != nil {
		updateData["Tanggal_Keluar"] = keluarPtr
	}

	log.Printf("[INACBG] Update data: %v\n", updateData)

	res := tx.Model(&models.BillingPasien{}).
		Where("ID_Billing = ?", input.ID_Billing).
		Updates(updateData)

	if res.Error != nil {
		tx.Rollback()
		log.Printf("[INACBG] Error updating billing: %v\n", res.Error)
		return fmt.Errorf("gagal update billing: %w", res.Error)
	}

	log.Printf("[INACBG] Updated %d rows in billing_pasien\n", res.RowsAffected)

	// 3. Bulk insert kode INACBG berdasarkan tipe_inacbg
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

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("[INACBG] Error committing transaction: %v\n", err)
		return err
	}

	log.Printf("[INACBG] ✅ Successfully saved INACBG for ID_Billing=%d, billing_sign=%s\n", input.ID_Billing, input.Billing_sign)

	// 4. Kirim email ke dokter jika billing_sign tidak kosong
	if input.Billing_sign != "" && strings.TrimSpace(input.Billing_sign) != "" {
		// Kirim email secara async (jika gagal, tidak mempengaruhi proses utama)
		// Log error jika ada, tapi tidak return error
		if err := SendEmailBillingSignToDokter(input.ID_Billing); err != nil {
			// Log error tapi tidak return error agar proses utama tetap berhasil
			// Di production, bisa menggunakan logger yang lebih proper
			fmt.Printf("Warning: Gagal mengirim email ke dokter untuk billing ID %d: %v\n", input.ID_Billing, err)
		}
	}

	return nil
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

	// Ambil INACBG RI
	inacbgRIMap := make(map[int][]string)
	var inacbgRIRows []struct {
		ID_Billing int
		Kode       string
	}
	if err := db.Table("billing_inacbg_ri").
		Where("ID_Billing IN ?", billingIDs).
		Select("ID_Billing, ID_INACBG_RI as Kode").
		Scan(&inacbgRIRows).Error; err != nil {
		return nil, err
	}
	for _, row := range inacbgRIRows {
		inacbgRIMap[row.ID_Billing] = append(inacbgRIMap[row.ID_Billing], row.Kode)
	}

	// Ambil INACBG RJ
	inacbgRJMap := make(map[int][]string)
	var inacbgRJRows []struct {
		ID_Billing int
		Kode       string
	}
	if err := db.Table("billing_inacbg_rj").
		Where("ID_Billing IN ?", billingIDs).
		Select("ID_Billing, ID_INACBG_RJ as Kode").
		Scan(&inacbgRJRows).Error; err != nil {
		return nil, err
	}
	for _, row := range inacbgRJRows {
		inacbgRJMap[row.ID_Billing] = append(inacbgRJMap[row.ID_Billing], row.Kode)
	}

	// Ambil dokter dari billing_dokter dengan urutan tanggal
	dokterMap := make(map[int][]string)
	var dokterRows []struct {
		ID_Billing int
		Nama       string
	}
	if err := db.Table("billing_dokter bd").
		Select("bd.ID_Billing, d.Nama_Dokter as Nama").
		Joins("JOIN dokter d ON bd.ID_Dokter = d.ID_Dokter").
		Where("bd.ID_Billing IN ?", billingIDs).
		Order("bd.Tanggal ASC").
		Scan(&dokterRows).Error; err != nil {
		return nil, err
	}
	for _, row := range dokterRows {
		dokterMap[row.ID_Billing] = append(dokterMap[row.ID_Billing], row.Nama)
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
			Total_Klaim:    b.Total_Tarif_BPJS,
			Tindakan_RS:    tindakanMap[b.ID_Billing],
			ICD9:           icd9Map[b.ID_Billing],
			ICD10:          icd10Map[b.ID_Billing],
			INACBG_RI:      inacbgRIMap[b.ID_Billing],
			INACBG_RJ:      inacbgRJMap[b.ID_Billing],
			Billing_sign:   b.Billing_sign,
			Nama_Dokter:    dokterMap[b.ID_Billing],
		}

		result = append(result, item)
	}

	return result, nil
}

// GetBillingByID - Get specific billing data by ID
func GetBillingByID(db *gorm.DB, id string) (map[string]interface{}, error) {
	var billing models.BillingPasien

	if err := db.Where("ID_Billing = ?", id).First(&billing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("billing dengan ID=%s tidak ditemukan", id)
		}
		return nil, fmt.Errorf("gagal mengambil billing: %w", err)
	}

	result := map[string]interface{}{
		"id_billing":     billing.ID_Billing,
		"id_pasien":      billing.ID_Pasien,
		"cara_bayar":     billing.Cara_Bayar,
		"tanggal_masuk":  billing.Tanggal_masuk,
		"tanggal_keluar": billing.Tanggal_keluar,
		"total_tarif_rs": billing.Total_Tarif_RS,
		"total_klaim":    billing.Total_Tarif_BPJS,
		"billing_sign":   billing.Billing_sign,
	}

	return result, nil
}
