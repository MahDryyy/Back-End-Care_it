package handlers

import (
	"net/http"

	"main.go/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", healthHandler)
	r.GET("/tarifBPJSRawatInap", listTarifBPJSRawatInapHandler)
	r.GET("/tarifBPJS/:kode", detailTarifBPJSRawatInapHandler)
	r.GET("/tarifBPJSRawatJalan", listTarifBPJSRawatJalanHandler)
	r.GET("/tarifBPJSRawatJalan/:kode", detailTarifBPJSRawatJalanHandler)
	r.GET("/tarifRS", listTarifRSHandler)
	r.GET("/tarifRS/:kode", detailTarifRSHandler)
	r.GET("/tarifRSByKategori/:kategori", listTarifRSByKategoriHandler)
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
