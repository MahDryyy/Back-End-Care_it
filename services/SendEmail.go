package services

import (
	"fmt"
	"net/smtp"
	"strings"

	"backendcareit/database"
	"backendcareit/models"

	"gorm.io/gorm"
)

// SendEmail mengirim email menggunakan SMTP
func SendEmail(to, subject, body string) error {
	from := "asikmahdi@gmail.com"
	password := "njom rhxb prrj tuoj"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("konfigurasi email tidak lengkap. Pastikan EMAIL_FROM, EMAIL_PASSWORD, SMTP_HOST, dan SMTP_PORT sudah di-set")
	}

	// Setup authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Format email message
	msg := []byte(fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	// Send email
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("gagal mengirim email: %w", err)
	}

	return nil
}

// SendEmailToMultiple mengirim email ke beberapa penerima sekaligus
func SendEmailToMultiple(to []string, subject, body string) error {
	from := "asikmahdi@gmail.com"
	password := "njom rhxb prrj tuoj"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("konfigurasi email tidak lengkap")
	}

	if len(to) == 0 {
		return fmt.Errorf("daftar penerima email tidak boleh kosong")
	}

	// Setup authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Format To header dengan semua penerima
	toHeader := strings.Join(to, ", ")

	// Format email message
	msg := []byte(fmt.Sprintf("To: %s\r\n", toHeader) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	// Send email ke semua penerima
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	err := smtp.SendMail(addr, auth, from, to, msg)
	if err != nil {
		return fmt.Errorf("gagal mengirim email: %w", err)
	}

	return nil
}

// SendEmailTest mengirim email test ke teman-teman
func SendEmailTest() error {
	to := []string{"stylohype685@gmail.com", "pasaribumonica2@gmail.com", "yestondehaan607@gmail.com"}
	subject := "Test Email - Sistem Billing Care IT"
	body := `
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
				.content { background-color: #f9f9f9; padding: 20px; margin-top: 20px; }
				.footer { margin-top: 20px; padding: 10px; text-align: center; font-size: 12px; color: #666; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h2>Test Email - Sistem Billing Care IT</h2>
				</div>
				<div class="content">
					<p>Halo!</p>
					<p>Ini adalah email test dari sistem billing Care IT.</p>
					<p>Jika Anda menerima email ini, berarti sistem email berfungsi dengan baik.</p>
					<p>Terima kasih!</p>
				</div>
				<div class="footer">
					<p>Sistem Billing Care IT</p>
					<p>Email ini dikirim untuk keperluan testing.</p>
				</div>
			</div>
		</body>
		</html>
	`

	if err := SendEmailToMultiple(to, subject, body); err != nil {
		return fmt.Errorf("gagal mengirim email test: %w", err)
	}

	return nil
}

// SendEmailBillingSignToDokter mengirim email ke semua dokter yang menangani pasien tentang billing sign
func SendEmailBillingSignToDokter(idBilling int) error {
	// 1. Ambil billing berdasarkan ID_Billing
	var billing models.BillingPasien
	if err := database.DB.Where("ID_Billing = ?", idBilling).First(&billing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("billing dengan ID_Billing=%d tidak ditemukan", idBilling)
		}
		return fmt.Errorf("gagal mengambil billing: %w", err)
	}

	// 2. Ambil semua dokter dari billing_dokter
	var dokterList []models.Dokter
	if err := database.DB.
		Table("billing_dokter bd").
		Select("d.*").
		Joins("JOIN dokter d ON bd.ID_Dokter = d.ID_Dokter").
		Where("bd.ID_Billing = ?", idBilling).
		Find(&dokterList).Error; err != nil {
		return fmt.Errorf("gagal mengambil dokter: %w", err)
	}

	if len(dokterList) == 0 {
		return fmt.Errorf("tidak ada dokter yang terkait dengan billing ID_Billing=%d", idBilling)
	}

	// 3. Ambil data pasien untuk informasi lengkap
	var pasien models.Pasien
	if err := database.DB.Where("ID_Pasien = ?", billing.ID_Pasien).First(&pasien).Error; err != nil {
		return fmt.Errorf("gagal mengambil data pasien: %w", err)
	}

	// 4. Format billing sign untuk ditampilkan
	billingSignDisplay := strings.ToUpper(billing.Billing_sign)
	if billingSignDisplay == "" {
		billingSignDisplay = "Belum ditentukan"
	}

	// 5. Kumpulkan semua email dokter yang valid
	var emailList []string
	var namaDokterList []string

	for _, dokter := range dokterList {
		// Pilih email dokter (prioritas Email_UB, jika kosong pakai Email_Pribadi)
		emailDokter := strings.TrimSpace(dokter.Email_UB)
		if emailDokter == "" {
			emailDokter = strings.TrimSpace(dokter.Email_Pribadi)
		}

		if emailDokter != "" {
			emailList = append(emailList, emailDokter)
			namaDokterList = append(namaDokterList, dokter.Nama_Dokter)
		}
	}

	if len(emailList) == 0 {
		return fmt.Errorf("tidak ada dokter dengan email yang terdaftar untuk billing ID_Billing=%d", idBilling)
	}

	// 6. Buat subject dan body email
	subject := fmt.Sprintf("Notifikasi Billing Sign - Pasien: %s", pasien.Nama_Pasien)

	// Gabungkan nama dokter untuk ditampilkan di email
	namaDokterStr := strings.Join(namaDokterList, ", ")

	body := fmt.Sprintf(`
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
				.content { background-color: #f9f9f9; padding: 20px; margin-top: 20px; }
				.info-row { margin: 10px 0; }
				.label { font-weight: bold; }
				.billing-sign { font-size: 18px; font-weight: bold; padding: 10px; text-align: center; margin: 20px 0; }
				.sign-hijau { background-color: #4CAF50; color: white; }
				.sign-kuning { background-color: #FFC107; color: #333; }
				.sign-orange { background-color: #FF9800; color: white; }
				.sign-merah { background-color: #F44336; color: white; }
				.footer { margin-top: 20px; padding: 10px; text-align: center; font-size: 12px; color: #666; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h2>Notifikasi Billing Sign</h2>
				</div>
				<div class="content">
					<p>Yth. Dr. %s,</p>
					<p>Berikut adalah informasi billing sign untuk pasien yang Anda tangani:</p>
					
					<div class="info-row">
						<span class="label">Nama Pasien:</span> %s
					</div>
					<div class="info-row">
						<span class="label">ID Billing:</span> %d
					</div>
					<div class="info-row">
						<span class="label">Ruangan:</span> %s
					</div>
					<div class="info-row">
						<span class="label">Kelas:</span> %s
					</div>
					<div class="info-row">
						<span class="label">Cara Bayar:</span> %s
					</div>
					<div class="info-row">
						<span class="label">Total Tarif RS:</span> Rp %.2f
					</div>
					<div class="info-row">
						<span class="label">Total Klaim BPJS:</span> Rp %.2f
					</div>
					
					<div class="billing-sign sign-%s">
						Billing Sign: %s
					</div>
					
					<p>Terima kasih atas perhatiannya.</p>
				</div>
				<div class="footer">
					<p>Sistem Billing Care IT</p>
					<p>Email ini dikirim secara otomatis, mohon tidak membalas email ini.</p>
				</div>
			</div>
		</body>
		</html>
	`, namaDokterStr, pasien.Nama_Pasien, billing.ID_Billing, pasien.Ruangan, pasien.Kelas,
		billing.Cara_Bayar, billing.Total_Tarif_RS, billing.Total_Tarif_BPJS,
		strings.ToLower(billing.Billing_sign), billingSignDisplay)

	// 7. Kirim email ke semua dokter
	if err := SendEmailToMultiple(emailList, subject, body); err != nil {
		return fmt.Errorf("gagal mengirim email ke dokter: %w", err)
	}

	return nil
}
