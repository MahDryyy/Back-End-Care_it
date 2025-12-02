package main

import (
	"fmt"
	"log"

	"backendcareit/database"
	"backendcareit/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.KonekDB()
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}
	database.DB = db

	r := gin.Default()
	r.Use(cors.Default())
	handlers.RegisterRoutes(r)

	port := "0.0.0.0:8081"
	fmt.Printf("Server berjalan di http://0.0.0.0:8081\n")
	fmt.Println("Akses dari jaringan lain menggunakan IP lokal komputer + port 8081")
	if err := r.Run(port); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
