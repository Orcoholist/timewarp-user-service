// main.go
package main

import (
	"log"
	"net/http"
	"os"

	"supabase-go-server/handlers"
	"supabase-go-server/middleware"
	"supabase-go-server/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.Timeline{},
		&models.TravelLog{},
		&models.User{},
		&models.TimeMachine{},
	)
}

func main() {
	// Загрузка .env
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден")
	}

	// Получение DATABASE_URL из .env
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("Переменная DATABASE_URL не задана в .env")
	}

	// Проверка существования файла сертификата
	certPath := os.Getenv("SSL_CERT_PATH")
	if certPath != "" {
		if _, err := os.Stat(certPath); os.IsNotExist(err) {
			log.Fatalf("Файл сертификата не найден: %s", certPath)
		}
		dbURL += "&sslrootcert=" + certPath
	}

	// Подключение к БД через GORM
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Автомиграция
	AutoMigrate(db)

	// Тестовое подключение к БД
	var version string
	if err := db.Raw("SELECT version()").Scan(&version).Error; err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	log.Println("Connected to:", version)

	// Создание нового мультиплексора
	mux := http.NewServeMux()

	// Регистрация маршрутов
	mux.HandleFunc("/api/timewarp", handlers.TimewarpHandler())
	mux.HandleFunc("/api/directions", handlers.GetDirections(db))
	mux.HandleFunc("/api/travel-logs", handlers.GetTravelLogs(db))
	mux.HandleFunc("/api/user", handlers.GetUserHandler(db))

	// Применение middleware к мультиплексору
	handler := middleware.LoggingMiddleware(mux)

	// Запуск сервера
	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
