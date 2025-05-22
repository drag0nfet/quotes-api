package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/drag0nfet/quotes-api/internal/models"
	"github.com/drag0nfet/quotes-api/internal/storage"
	"github.com/gorilla/mux"
)

// Handler - обёртка для хранения переменной хранилища при выполнении запросов
type Handler struct {
	store *storage.MemoryStorage
}

func NewHandler(store *storage.MemoryStorage) *Handler {
	return &Handler{store: store}
}

// AddQuote - POST /quotes
func (h *Handler) AddQuote(w http.ResponseWriter, r *http.Request) {
	var quote models.Quote

	if err := json.NewDecoder(r.Body).Decode(&quote); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	if quote.Author == "" || quote.Text == "" {
		http.Error(w, "Автор и цитата обязательны", http.StatusBadRequest)
		return
	}

	id := h.store.AddQuote(quote)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// GetQuotes - GET /quotes
func (h *Handler) GetQuotes(w http.ResponseWriter, r *http.Request) {
	quotes := h.store.GetAllQuotes()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}

// GetRandomQuote - GET /quotes/random
func (h *Handler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, exists := h.store.GetRandomQuote()
	if !exists {
		http.Error(w, "Нет добавленных цитат", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

// GetQuotesByAuthor - GET /quotes?author=<name>
func (h *Handler) GetQuotesByAuthor(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	if author == "" {
		h.GetQuotes(w, r)
		return
	}

	quotes := h.store.GetQuotesByAuthor(author)
	if len(quotes) == 0 {
		http.Error(w, "Нет цитат заданного автора", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}

// DeleteQuote - DELETE /quotes/{id}
func (h *Handler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	if !h.store.DeleteQuote(id) {
		http.Error(w, "Цитата не найдена", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
