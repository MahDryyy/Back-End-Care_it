package handlers

import (
	"net/http"
	"strconv"

	"main.go/models"
	"main.go/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/dokter", listDokterHandler)
	r.GET("/ruangan", listRuanganHandler)
	r.GET("/icd10", listICD10Handler)
	r.GET("/icd9", listICD9Handler)
	r.GET("/", healthHandler)
	r.GET("/tarifBPJSRawatInap", listTarifBPJSRawatInapHandler)
	r.GET("/tarifBPJS/:kode", detailTarifBPJSRawatInapHandler)
	r.GET("/tarifBPJSRawatJalan", listTarifBPJSRawatJalanHandler)
	r.GET("/tarifBPJSRawatJalan/:kode", detailTarifBPJSRawatJalanHandler)
	r.GET("/tarifRS", listTarifRSHandler)
	r.GET("/tarifRS/:kode", detailTarifRSHandler)
	r.GET("/tarifRSByKategori/:kategori", listTarifRSByKategoriHandler)
	r.GET("/pasien/:id", GetPasien)
	r.POST("/billing", CreateBillingHandler)
}

// Health check
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Server berjalan",
	})
}

// List tarif BPJS Rawat Inap
func listTarifBPJSRawatInapHandler(c *gin.Context) {
	data, err := services.GetTarifBPJSRawatInap()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data",
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func detailTarifBPJSRawatInapHandler(c *gin.Context) {
	kode := c.Param("kode")
	data, err := services.GetTarifBPJSRawatInapByKode(kode)
	if err != nil {
		if services.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "not_found",
				"message": "Kode tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data",
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

// List tarif BPJS Rawat Jalan
func listTarifBPJSRawatJalanHandler(c *gin.Context) {
	data, err := services.GetTarifBPJSRawatJalan()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data",
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func detailTarifBPJSRawatJalanHandler(c *gin.Context) {
	kode := c.Param("kode")
	data, err := services.GetTarifBPJSRawatJalanByKode(kode)
	if err != nil {
		if services.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "not_found",
				"message": "Kode tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data",
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

// List tarif RS
func listTarifRSHandler(c *gin.Context) {
	data, err := services.GetTarifRS()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data",
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func detailTarifRSHandler(c *gin.Context) {
	kode := c.Param("kode")
	data, err := services.GetTarifRSByKode(kode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data",
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func listTarifRSByKategoriHandler(c *gin.Context) {
	kategori := c.Param("kategori")
	data, err := services.GetTarifRSByKategori(kategori)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data",
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

// ICD9
func listICD9Handler(c *gin.Context) {
	data, err := services.GetICD9()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data",
		})
		return
	}
	c.JSON(http.StatusOK, data)
}

// ICD10
func listICD10Handler(c *gin.Context) {
	data, err := services.GetICD10()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data",
		})
		return
	}
	c.JSON(http.StatusOK, data)
}

// ruangan
func listRuanganHandler(c *gin.Context) {
	data, err := services.GetRuangan()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data",
		})
		return
	}
	c.JSON(http.StatusOK, data)
}

// dokter
func listDokterHandler(c *gin.Context) {
	data, err := services.GetDokter()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data",
		})
		return
	}
	c.JSON(http.StatusOK, data)
}

//Liat pasien sudah atau belum

func GetPasien(c *gin.Context) {
	idStr := c.Param("id")

	// Konversi string ke int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "ID pasien harus berupa angka",
		})
		return
	}

	pasien, err := services.GetPasienByID(id)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Pasien tidak ditemukan",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Data pasien ditemukan",
		"data":    pasien,
	})
}

//add pasien baru

// CreateBillingHandler handler untuk membuat billing baru dari data frontend
func CreateBillingHandler(c *gin.Context) {
	var input models.BillingRequest

	// Pastikan JSON
	contentType := c.GetHeader("Content-Type")
	if contentType != "application/json" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Content-Type harus application/json",
			"error":   "Content-Type yang diterima: " + contentType,
		})
		return
	}

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Panggil service → return 5 data
	billing, pasien, tindakanList, icd9List, icd10List, err :=
		services.DataFromFE(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat billing",
			"error":   err.Error(),
		})
		return
	}

	// Response lengkap ke FE
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Billing berhasil dibuat",
		"data": gin.H{
			"pasien":      pasien,
			"billing":     billing,
			"tindakan_rs": tindakanList,
			"icd9":        icd9List,
			"icd10":       icd10List,
		},
	})
}

