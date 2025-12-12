package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"backendcareit/models"
)

// loginRequest merepresentasikan payload login dokter.

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}
	return "SECRET_KAMU"
}

// LoginDokterHandler membuat handler POST /login yang memverifikasi kredensial
// dokter dan mengembalikan JWT jika sukses.

// GenerateToken membuat JWT yang berlaku 24 jam untuk dokter.
func GenerateToken(dokter models.Dokter, email string) (string, error) {
	claims := jwt.MapClaims{
		"id":    dokter.ID_Dokter,
		"nama":  dokter.Nama_Dokter,
		"ksm":   dokter.KSM,
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateTokenAdmin membuat JWT yang berlaku 24 jam untuk admin.
func GenerateTokenAdmin(admin models.Admin_Ruangan, namaAdmin string) (string, error) {
	claims := jwt.MapClaims{
		"id":         admin.ID_Admin,
		"nama_admin": admin.Nama_Admin,
		"id_ruangan": admin.ID_Ruangan,
		"role":       "admin",
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// AuthMiddleware memvalidasi header Authorization Bearer <token> dan
// menyimpan data dokter pada context jika token valid.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Authorization header wajib menggunakan Bearer token",
			})
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("metode tanda tangan tidak dikenal")
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Token tidak valid atau kadaluarsa",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Token tidak valid",
			})
			return
		}

		c.Set("dokter_id", claims["id"])
		c.Set("dokter_nama", claims["nama"])
		c.Set("dokter_ksm", claims["ksm"])
		c.Set("dokter_email", claims["email"])

		c.Next()
	}
}
