package main

import (
	"fmt"
	"log"

	"main.go/database"
	"main.go/handlers"

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

	port := ":8081"
	fmt.Printf("Server berjalan di http://localhost%s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
