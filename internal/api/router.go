package api

import (
	"github.com/drag0nfet/quotes-api/internal/storage"
	"github.com/gorilla/mux"
)

// SetupRoutes настраивает доступные в API маршруты
func SetupRoutes(router *mux.Router, store *storage.MemoryStorage) {
	h := NewHandler(store)

	router.HandleFunc("/quotes", h.AddQuote).Methods("POST")
	router.HandleFunc("/quotes", h.GetQuotesByAuthor).Methods("GET")
	router.HandleFunc("/quotes/random", h.GetRandomQuote).Methods("GET")
	router.HandleFunc("/quotes/{id}", h.DeleteQuoteHandler).Methods("DELETE")
}
