package services

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"main.go/database"
	"main.go/models"
)

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

func DataFromFE(input models.BillingRequest) (*models.BillingPasien, *models.Pasien, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Cek apakah pasien dengan nama tersebut sudah ada, jika tidak buat baru
	var pasien models.Pasien
	result := tx.Where("Nama_Pasien = ?", input.Nama_Pasien).First(&pasien)

	if result.Error != nil {
		// Jika pasien tidak ditemukan, buat pasien baru
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Buat pasien baru (ID_Pasien akan auto-increment oleh database)
			pasien = models.Pasien{
				Nama_Pasien:   input.Nama_Pasien,
				Jenis_Kelamin: input.Jenis_Kelamin,
				Usia:          input.Usia,
				Ruangan:       input.Ruangan,
				Kelas:         input.Kelas,
			}

			// Create pasien dalam transaction (ID_Pasien akan di-generate oleh database)
			if err := tx.Create(&pasien).Error; err != nil {
				tx.Rollback()
				return nil, nil, errors.New("gagal membuat pasien baru: " + err.Error())
			}

			// Setelah Create, pasien sudah terisi dengan data yang benar termasuk ID_Pasien (auto-increment)
		} else {
			// Error lain selain record not found
			tx.Rollback()
			return nil, nil, errors.New("gagal mencari pasien: " + result.Error.Error())
		}
	}
	// Jika pasien ditemukan, gunakan data yang ada (pasien sudah terisi dari query di atas)

	// Pastikan ID_Pasien valid sebelum membuat billing
	if pasien.ID_Pasien == 0 {
		tx.Rollback()
		return nil, nil, errors.New("ID_Pasien tidak valid")
	}

	// Cari ID_Dokter berdasarkan Nama_Dokter
	var dokter models.Dokter
	if err := tx.Where("Nama_Dokter = ?", input.Nama_Dokter).First(&dokter).Error; err != nil {
		tx.Rollback()
		return nil, nil, errors.New("dokter dengan nama " + input.Nama_Dokter + " tidak ditemukan")
	}

	// Buat billing baru (ID_Billing akan auto-increment oleh database)
	now := time.Now()
	billing := models.BillingPasien{
		ID_Pasien:        pasien.ID_Pasien, // Gunakan ID_Pasien dari pasien yang ditemukan/dibuat
		Cara_Bayar:       input.Cara_Bayar,
		Tanggal_masuk:    &now,
		ID_Dokter:        dokter.ID_Dokter, // Gunakan ID_Dokter dari dokter yang ditemukan
		Total_Tarif_RS:   input.Total_Tarif_RS,
		Total_Tarif_BPJS: 0, // Default 0, bisa dihitung nanti jika diperlukan
		Billing_sign:     "created",
	}

	if err := tx.Create(&billing).Error; err != nil {
		tx.Rollback()
		return nil, nil, errors.New("gagal membuat billing: " + err.Error() + " | ID_Pasien yang digunakan: " + fmt.Sprintf("%d", pasien.ID_Pasien))
	}

	if err := tx.Commit().Error; err != nil {
		return nil, nil, err
	}

	return &billing, &pasien, nil
}
