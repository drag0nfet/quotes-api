package storage

import (
	"github.com/drag0nfet/quotes-api/internal/models"
	"sync"
)

type MemoryStorage struct {
	quotes      map[int]models.Quote // Хранение цитат. 			Ключ - id, 		Значение - цитата
	authorIndex map[string][]int     // Хранение цитат автора. 	Ключ - автор, 	Значение - слайс id цитат
	nextID      int                  // Новая цитата будет иметь этот id.
	mutex       sync.RWMutex         // Для корректного конкурентного доступа к чтению/записи.
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		quotes:      make(map[int]models.Quote),
		authorIndex: make(map[string][]int),
		nextID:      1,
	}
}

func (s *MemoryStorage) AddQuote(quote models.Quote) int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	quote.ID = s.nextID
	s.quotes[quote.ID] = quote
	s.authorIndex[quote.Author] = append(s.authorIndex[quote.Author], quote.ID)
	s.nextID++

	return quote.ID
}

func (s *MemoryStorage) GetAllQuotes() []models.Quote {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	quotes := make([]models.Quote, 0, len(s.quotes))
	for _, quote := range s.quotes {
		quotes = append(quotes, quote)
	}
	return quotes
}

func (s *MemoryStorage) GetRandomQuote() (models.Quote, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if len(s.quotes) == 0 {
		return models.Quote{}, false
	}

	var randomID int
	for key := range s.quotes {
		randomID = key
		break
	}

	quote, exists := s.quotes[randomID]
	return quote, exists
}

func (s *MemoryStorage) GetQuotesByAuthor(author string) []models.Quote {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	quoteIDs, exists := s.authorIndex[author]
	if !exists {
		return nil
	}

	quotes := make([]models.Quote, 0, len(quoteIDs))
	for _, id := range quoteIDs {
		if quote, ok := s.quotes[id]; ok {
			quotes = append(quotes, quote)
		}
	}
	return quotes
}

// DeleteQuote удаляет цитату по ID.
func (s *MemoryStorage) DeleteQuote(id int) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	quote, exists := s.quotes[id]
	if !exists {
		return false
	}

	delete(s.quotes, id)

	authorQuotes := s.authorIndex[quote.Author]
	for i, quoteID := range authorQuotes {
		if quoteID == id {
			s.authorIndex[quote.Author] = append(authorQuotes[:i], authorQuotes[i+1:]...)
			break
		}
	}

	if len(s.authorIndex[quote.Author]) == 0 {
		delete(s.authorIndex, quote.Author)
	}

	return true
}
