package main

import (
	"log"
	"net/http"

	"github.com/drag0nfet/quotes-api/internal/api"
	"github.com/drag0nfet/quotes-api/internal/storage"
	"github.com/gorilla/mux"
)

func main() {
	// Определение временного хранилища для цитат на время работы сервера
	store := storage.NewMemoryStorage()

	// Настройка маршрутов
	router := mux.NewRouter()
	api.SetupRoutes(router, store)

	// Запуск сервера
	const port = ":8080"
	log.Printf("Сервер запущен на http://localhost%s", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
