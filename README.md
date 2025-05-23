# Quotes API

REST API сервис для управления цитатами, реализованный на Go. Поддерживает добавление, получение, фильтрацию и удаление цитат.

## Требования

- Go 1.21 или выше
- Установленный пакет `gorilla/mux` (`go get github.com/gorilla/mux`)

## Установка

1. Склонируйте репозиторий:
   `git clone https://github.com/drag0nfet/quotes-api.git`
   `cd quotes-api`

2. Установите зависимости:
   `go mod tidy`

3. Запустите приложение:
   `go run ./cmd/quotes-api`

Сервер будет доступен по адресу `http://localhost:8080`.

## Использование

API поддерживает следующие эндпоинты:

- **POST /quotes** — добавление новой цитаты. Ожидаемый ответ: JSON вида `{"id":id}`.
- **GET /quotes** — получение всех цитат. Ожидаемый ответ: массив цитат в формате JSON вида `[{"id":id,"author":author,"quote":quote},...]` или пустой JSON-массив при отсутствии цитат в хранилище.
- **GET /quotes/random** — получение случайной цитаты. Ожидаемый ответ: цитата в формате JSON вида `{"id":id,"author":author,"quote":quote}` или сообщение `Нет добавленных цитат` (HTTP-404).
- **GET /quotes?author={Author}** — получение цитат по автору. Ожидаемый ответ совпадает с **GET /quotes**, но при отсутствии запрашиваемых цитат поступит сообщение `Нет цитат заданного автора` (HTTP-404).
- **DELETE /quotes/{id}** — удаление цитаты по ID. Ожидаемый ответ: в случае успеха - StatusNoContent (HTTP-204), при отсутствии цитаты с заданным id - сообщение `Цитата не найдена` (HTTP-404).

### Примеры консольных запросов

1. Добавление цитаты:
   `curl -X POST http://localhost:8080/quotes -H "Content-Type: application/json" -d '{"author":"Confucius","quote":"Life is simple, but we insist on making it complicated."}'`
 
2. Получение всех цитат:
   `curl http://localhost:8080/quotes`

3. Получение случайной цитаты:
   `curl http://localhost:8080/quotes/random`

4. Фильтрация цитат по автору:
   `curl http://localhost:8080/quotes?author=Confucius`

5. Удаление цитаты по ID:
   `curl -X DELETE http://localhost:8080/quotes/1`

## Тестирование

Для запуска юнит-тестов, проверяющих обработчики API и взаимодействие с хранилищем цитат, выполните: `go test ./internal/tests -v`

Тесты покрывают следующие аспекты:

- Корректность HTTP-ответов для всех эндпоинтов
- Поведение in-memory хранилища (добавление, получение и удаление цитат)
- Пограничные случаи:
   - Некорректный JSON
   - Пустые поля
   - Несуществующие ID и авторы

Для проверки покрытия кода используйте: `go test -coverpkg=./... ./internal/tests`

## Структура проекта

- `cmd/quotes-api/main.go` — точка входа приложения.
- `internal/api/` — обработчики HTTP-запросов и маршрутизация.
- `internal/models/` — структура данных Quote.
- `internal/storage/` — логика хранения цитат в памяти.
- `internal/tests/` — юнит-тесты.

## Заметки

- Данные хранятся в памяти и не сохраняются между запусками сервера.
- При запуске сервера инициализируется пустое хранилище, первая добавленная цитата будет иметь id=1.
- При удалении цитаты, её id не будет использоваться до перезапуска сервера.
