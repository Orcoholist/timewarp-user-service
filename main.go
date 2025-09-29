package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// Модели данных
type Timewarp struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Description string    `json:"description"`
	Count     *int64    `json:"count,omitempty"`
}

type Direction struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Year        int    `json:"year"`
	Description *string `json:"description,omitempty"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

// Функции работы с БД
func createTimewarp(pool *pgxpool.Pool, description string, count *int64) (int64, error) {
	query := `
		INSERT INTO timewarp (description, count)
		VALUES ($1, $2)
		RETURNING id`
	var id int64
	err := pool.QueryRow(context.Background(), query, description, count).Scan(&id)
	return id, err
}

func getDirections(pool *pgxpool.Pool) ([]Direction, error) {
	rows, err := pool.Query(context.Background(), "SELECT id, name, year, description FROM direction")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var directions []Direction
	for rows.Next() {
		var d Direction
		err := rows.Scan(&d.ID, &d.Name, &d.Year, &d.Description)
		if err != nil {
			return nil, err
		}
		directions = append(directions, d)
	}
	return directions, nil
}

func getUserByUsername(pool *pgxpool.Pool, username string) (*User, error) {
	var u User
	err := pool.QueryRow(context.Background(), "SELECT id, username, password FROM users WHERE username = $1", username).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// HTTP-обработчики
func handleCreateTimewarp(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Поступил запрос: %s %s", r.Method, r.URL.Path) // Логируем запрос

		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			log.Printf("Ошибка: Неверный метод %s", r.Method)
			return
		}

		var req struct {
			Description string `json:"description"`
			Count       *int64 `json:"count,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
			log.Printf("Ошибка декодирования JSON: %v", err)
			return
		}

		id, err := createTimewarp(pool, req.Description, req.Count)
		if err != nil {
			http.Error(w, "Ошибка при создании записи", http.StatusInternalServerError)
			log.Printf("Ошибка при создании timewarp: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int64{"id": id})
		log.Printf("Успешный ответ: ID %d", id)
	}
}

func handleGetDirections(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		directions, err := getDirections(pool)
		if err != nil {
			http.Error(w, "Ошибка при получении направлений", http.StatusInternalServerError)
			log.Println("Ошибка:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(directions)
	}
}

func handleGetUser(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		username := r.URL.Query().Get("username")
		if username == "" {
			http.Error(w, "Не указан username", http.StatusBadRequest)
			return
		}

		user, err := getUserByUsername(pool, username)
		if err != nil {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func main() {
	// Загрузка .env
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден")
	}

	// Получение переменных окружения
	dbURL := os.Getenv("DATABASE_URL")
	certPath := os.Getenv("SSL_CERT_PATH")

	// Проверка существования файла сертификата
	if certPath != "" {
		if _, err := os.Stat(certPath); os.IsNotExist(err) {
			log.Fatalf("Файл сертификата не найден: %s", certPath)
		}
	}

	// Добавление параметра sslrootcert в URL
	dbURL += "&sslrootcert=" + certPath

	// Создание пула подключений
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	// Тестовое подключение к БД
	var version string
	if err := pool.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	log.Println("Connected to:", version)

	// Регистрация маршрутов
	http.HandleFunc("/api/timewarp", handleCreateTimewarp(pool))
	http.HandleFunc("/api/directions", handleGetDirections(pool))
	http.HandleFunc("/api/user", handleGetUser(pool))

	// Запуск сервера
	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}