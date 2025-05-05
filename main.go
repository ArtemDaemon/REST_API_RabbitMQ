package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"rest-api/handlers"
	"rest-api/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	_ "rest-api/docs" // автосгенерированная документация

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Пример REST API
// @version 1.0
// @description Это REST API с авторизацией через Bearer-токен
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Aithorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
	}

	expectedToken := os.Getenv("API_TOKEN")
	if expectedToken == "" {
		log.Fatal("API_TOKEN не задан")
	}

	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal("Ошибка открытия базы данных:", err)
	}
	defer db.Close()

	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(expectedToken))

		r.Get("/count", handlers.CountItemsHandler(db))
		r.Get("/last_created_at", handlers.LastCreatedAtHandler(db))
		r.Get("/get_item", handlers.GetItemByDateHandler(db))
		r.Post("/add_item", handlers.AddItemHandler(db))
	})

	fmt.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", r)
}
