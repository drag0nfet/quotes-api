package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drag0nfet/quotes-api/internal/api"
	"github.com/drag0nfet/quotes-api/internal/models"
	"github.com/drag0nfet/quotes-api/internal/storage"
	"github.com/gorilla/mux"
)

// TestAddQuoteHandler тестирует POST /quotes.
func TestAddQuoteHandler(t *testing.T) {
	store := storage.NewMemoryStorage()
	handler := api.NewHandler(store)

	tests := []struct {
		name           string
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "ValidInput",
			body:           `{"author":"Confucius","quote":"Life is simple."}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   "{\"id\":1}\n",
		},
		{
			name:           "InvalidJSON",
			body:           `{"author":"Confucius",}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Некорректный JSON\n",
		},
		{
			name:           "EmptyAuthor",
			body:           `{"author":"","quote":"Life is simple."}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Автор и цитата обязательны\n",
		},
		{
			name:           "EmptyQuote",
			body:           `{"author":"Confucius","quote":""}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Автор и цитата обязательны\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/quotes", bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.AddQuote(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
			if rr.Body.String() != tt.expectedBody {
				t.Errorf("Expected body %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

// TestGetQuotesHandler тестирует GET /quotes.
func TestGetQuotesHandler(t *testing.T) {
	store := storage.NewMemoryStorage()
	handler := api.NewHandler(store)

	store.AddQuote(models.Quote{Author: "Confucius", Text: "Life is simple."})

	req, err := http.NewRequest("GET", "/quotes", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.GetQuotes(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var quotes []models.Quote
	if err := json.NewDecoder(rr.Body).Decode(&quotes); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(quotes) != 1 || quotes[0].Author != "Confucius" || quotes[0].Text != "Life is simple." {
		t.Errorf("Unexpected response: %v", quotes)
	}
}

// TestGetRandomQuoteHandler тестирует GET /quotes/random.
func TestGetRandomQuoteHandler(t *testing.T) {
	store := storage.NewMemoryStorage()
	handler := api.NewHandler(store)

	req, err := http.NewRequest("GET", "/quotes/random", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.GetRandomQuote(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
	if rr.Body.String() != "Нет добавленных цитат\n" {
		t.Errorf("Expected body %q, got %q", "Нет добавленных цитат\n", rr.Body.String())
	}

	store.AddQuote(models.Quote{Author: "Confucius", Text: "Life is simple."})

	rr = httptest.NewRecorder()
	handler.GetRandomQuote(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var quote models.Quote
	if err := json.NewDecoder(rr.Body).Decode(&quote); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if quote.Author != "Confucius" || quote.Text != "Life is simple." {
		t.Errorf("Unexpected response: %v", quote)
	}
}

// TestGetQuotesByAuthorHandler тестирует GET /quotes?author=<name>.
func TestGetQuotesByAuthorHandler(t *testing.T) {
	store := storage.NewMemoryStorage()
	handler := api.NewHandler(store)

	// Добавляем тестовые цитаты
	store.AddQuote(models.Quote{Author: "Confucius", Text: "Life is simple."})
	store.AddQuote(models.Quote{Author: "Einstein", Text: "Imagination is everything."})

	tests := []struct {
		name           string
		author         string
		expectedStatus int
		expectedQuotes int
	}{
		{
			name:           "ExistingAuthor",
			author:         "Confucius",
			expectedStatus: http.StatusOK,
			expectedQuotes: 1,
		},
		{
			name:           "NonExistingAuthor",
			author:         "Unknown",
			expectedStatus: http.StatusNotFound,
			expectedQuotes: 0,
		},
		{
			name:           "EmptyAuthor",
			author:         "",
			expectedStatus: http.StatusOK,
			expectedQuotes: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/quotes?author="+tt.author, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler.GetQuotesByAuthor(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			var quotes []models.Quote
			if err := json.NewDecoder(rr.Body).Decode(&quotes); err != nil && tt.expectedStatus == http.StatusOK {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if len(quotes) != tt.expectedQuotes {
				t.Errorf("Expected %d quotes, got %d", tt.expectedQuotes, len(quotes))
			}
		})
	}
}

// TestDeleteQuoteHandler тестирует DELETE /quotes/{id}.
func TestDeleteQuoteHandler(t *testing.T) {
	store := storage.NewMemoryStorage()
	handler := api.NewHandler(store)

	// Добавляем тестовую цитату
	store.AddQuote(models.Quote{Author: "Confucius", Text: "Life is simple."})

	tests := []struct {
		name           string
		id             string
		expectedStatus int
	}{
		{
			name:           "ValidID",
			id:             "1",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "InvalidID",
			id:             "invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "NonExistingID",
			id:             "999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/quotes/"+tt.id, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/quotes/{id}", handler.DeleteQuote).Methods("DELETE")
			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
