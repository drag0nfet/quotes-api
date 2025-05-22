package tests

import (
	"testing"

	"github.com/drag0nfet/quotes-api/internal/models"
	"github.com/drag0nfet/quotes-api/internal/storage"
)

func TestAddQuoteStorage(t *testing.T) {
	store := storage.NewMemoryStorage()

	tests := []struct {
		name                 string
		quote                models.Quote
		expectedID           int
		expectedLen          int
		expectedAuthorQuotes int
	}{
		{
			name:                 "ValidQuote",
			quote:                models.Quote{Author: "Confucius", Text: "Life is simple."},
			expectedID:           1,
			expectedLen:          1,
			expectedAuthorQuotes: 1,
		},
		{
			name:                 "EmptyAuthor",
			quote:                models.Quote{Author: "", Text: "No author."},
			expectedID:           1,
			expectedLen:          1,
			expectedAuthorQuotes: 1,
		},
		{
			name:                 "EmptyText",
			quote:                models.Quote{Author: "Einstein", Text: ""},
			expectedID:           1,
			expectedLen:          1,
			expectedAuthorQuotes: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := store.AddQuote(tt.quote)
			if id != tt.expectedID {
				t.Errorf("Expected ID %d, got %d", tt.expectedID, id)
			}

			quotes := store.GetAllQuotes()
			if len(quotes) != tt.expectedLen {
				t.Errorf("Expected %d quotes, got %d", tt.expectedLen, len(quotes))
			}

			if len(quotes) > 0 && (quotes[0].Author != tt.quote.Author || quotes[0].Text != tt.quote.Text) {
				t.Errorf("Unexpected quote: got %v, want %v", quotes[0], tt.quote)
			}

			authorQuotes := store.GetQuotesByAuthor(tt.quote.Author)
			if len(authorQuotes) != tt.expectedAuthorQuotes {
				t.Errorf("Expected %d quotes for author %q, got %d", tt.expectedAuthorQuotes, tt.quote.Author, len(authorQuotes))
			}

			// Сбрасываем хранилище для следующего теста
			store = storage.NewMemoryStorage()
		})
	}
}

func TestGetAllQuotesStorage(t *testing.T) {
	store := storage.NewMemoryStorage()

	tests := []struct {
		name        string
		addQuotes   []models.Quote
		expectedLen int
	}{
		{
			name:        "EmptyStore",
			addQuotes:   []models.Quote{},
			expectedLen: 0,
		},
		{
			name:        "SingleQuote",
			addQuotes:   []models.Quote{{Author: "Confucius", Text: "Life is simple."}},
			expectedLen: 1,
		},
		{
			name: "MultipleQuotes",
			addQuotes: []models.Quote{
				{Author: "Confucius", Text: "Life is simple."},
				{Author: "Einstein", Text: "Imagination is everything."},
			},
			expectedLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store = storage.NewMemoryStorage()
			for _, quote := range tt.addQuotes {
				store.AddQuote(quote)
			}

			quotes := store.GetAllQuotes()
			if len(quotes) != tt.expectedLen {
				t.Errorf("Expected %d quotes, got %d", tt.expectedLen, len(quotes))
			}

			if tt.expectedLen > 0 {
				for i, quote := range quotes {
					if i < len(tt.addQuotes) && (quote.Author != tt.addQuotes[i].Author || quote.Text != tt.addQuotes[i].Text) {
						t.Errorf("Unexpected quote at index %d: got %v, want %v", i, quote, tt.addQuotes[i])
					}
				}
			}
		})
	}
}

func TestGetRandomQuoteStorage(t *testing.T) {
	store := storage.NewMemoryStorage()

	tests := []struct {
		name           string
		addQuote       bool
		quote          models.Quote
		expectedExists bool
	}{
		{
			name:           "EmptyStore",
			addQuote:       false,
			expectedExists: false,
		},
		{
			name:           "SingleQuote",
			addQuote:       true,
			quote:          models.Quote{Author: "Confucius", Text: "Life is simple."},
			expectedExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store = storage.NewMemoryStorage()
			if tt.addQuote {
				store.AddQuote(tt.quote)
			}

			quote, exists := store.GetRandomQuote()
			if exists != tt.expectedExists {
				t.Errorf("Expected exists %v, got %v", tt.expectedExists, exists)
			}

			if tt.addQuote && (quote.Author != tt.quote.Author || quote.Text != tt.quote.Text) {
				t.Errorf("Unexpected quote: got %v, want %v", quote, tt.quote)
			}
		})
	}
}

func TestGetQuotesByAuthorStorage(t *testing.T) {
	store := storage.NewMemoryStorage()

	store.AddQuote(models.Quote{Author: "Confucius", Text: "Life is simple."})
	store.AddQuote(models.Quote{Author: "Confucius", Text: "Silence is a true friend."})
	store.AddQuote(models.Quote{Author: "Einstein", Text: "Imagination is everything."})

	tests := []struct {
		name           string
		author         string
		expectedQuotes int
	}{
		{
			name:           "ExistingAuthor",
			author:         "Confucius",
			expectedQuotes: 2,
		},
		{
			name:           "NonExistingAuthor",
			author:         "Unknown",
			expectedQuotes: 0,
		},
		{
			name:           "EmptyAuthor",
			author:         "",
			expectedQuotes: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quotes := store.GetQuotesByAuthor(tt.author)
			if len(quotes) != tt.expectedQuotes {
				t.Errorf("Expected %d quotes, got %d", tt.expectedQuotes, len(quotes))
			}

			if tt.author == "Confucius" && tt.expectedQuotes > 0 {
				for _, quote := range quotes {
					if quote.Author != "Confucius" {
						t.Errorf("Unexpected author: got %q, want %q", quote.Author, "Confucius")
					}
				}
			}
		})
	}
}

func TestDeleteQuoteStorage(t *testing.T) {
	store := storage.NewMemoryStorage()

	store.AddQuote(models.Quote{Author: "Confucius", Text: "Life is simple."})

	tests := []struct {
		name                 string
		id                   int
		expectedResult       bool
		expectedLen          int
		expectedAuthorQuotes int
	}{
		{
			name:                 "ValidID",
			id:                   1,
			expectedResult:       true,
			expectedLen:          0,
			expectedAuthorQuotes: 0,
		},
		{
			name:                 "NonExistingID",
			id:                   999,
			expectedResult:       false,
			expectedLen:          1,
			expectedAuthorQuotes: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store = storage.NewMemoryStorage()
			store.AddQuote(models.Quote{Author: "Confucius", Text: "Life is simple."})

			result := store.DeleteQuote(tt.id)
			if result != tt.expectedResult {
				t.Errorf("Expected result %v, got %v", tt.expectedResult, result)
			}

			quotes := store.GetAllQuotes()
			if len(quotes) != tt.expectedLen {
				t.Errorf("Expected %d quotes, got %d", tt.expectedLen, len(quotes))
			}

			authorQuotes := store.GetQuotesByAuthor("Confucius")
			if len(authorQuotes) != tt.expectedAuthorQuotes {
				t.Errorf("Expected %d quotes for author, got %d", tt.expectedAuthorQuotes, len(authorQuotes))
			}
		})
	}
}
