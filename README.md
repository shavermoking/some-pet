# some-pet

## Запуск

1. Клонируем репозиторий

```bash
git clone https://github.com/shavermoking/some-pet.git
cd some-pet
```

2. Создаём .env

```bash
cp .env.example .env
```

3. Запуск с Docker

```bash
docker-compose up -d
```
Сервис будет доступен на http://localhost:8080


## Эндпоинты

### 1. Создание книги
`POST http://localhost:8080/books`


```json
{
    "title": "Война и мир",
    "author": "Лев Толстой",
    "year": 1869,
    "isbn": "978-5-17-089835-2",
    "rating": 5
}
```

#### Response
```json
{
    "id": 14,
    "title": "Война и мир",
    "author": "Лев Толстой",
    "year": 1869,
    "isbn": "978-5-17-089835-2",
    "outOfStock": false,
    "rating": 5,
    "CreatedAt": "0001-01-01T00:00:00Z",
    "UpdatedAt": "0001-01-01T00:00:00Z"
}
```

### 2. Получение всех книг
`GET http://localhost:8080/books`


#### Response
```json
[
    {
        "id": 1,
        "title": "Война и мир",
        "author": "Лев Толстой",
        "year": 1869,
        "isbn": "978-5-17-089835-2",
        "rating": 5,
        "out_of_stock": false,
        "CreatedAt": "0001-01-01T00:00:00Z",
        "UpdatedAt": "0001-01-01T00:00:00Z"
    },
    {
        "id": 2,
        "title": "Преступление и наказание",
        "author": "Федор Достоевский",
        "year": 1866,
        "isbn": "978-5-04-096745-4",
        "rating": 5,
        "out_of_stock": true,
        "CreatedAt": "0001-01-01T00:00:00Z",
        "UpdatedAt": "0001-01-01T00:00:00Z"
    }
]
```

### 3. Получение книги по ID
`GET http://localhost:8080/books/1`

#### Response
```json
{
    "id": 1,
    "title": "Война и мир",
    "author": "Лев Толстой",
    "year": 1869,
    "isbn": "978-5-17-089835-2",
    "rating": 5,
    "out_of_stock": false,
    "CreatedAt": "0001-01-01T00:00:00Z",
    "UpdatedAt": "0001-01-01T00:00:00Z"
}
```

### 4. Обновление книги
`PUT http://localhost:8080/books/1`

```json
{
    "title": "Война и мир (обновленное издание)",
    "author": "Лев Толстой",
    "year": 1873,
    "isbn": "978-5-17-089835-2",
    "rating": 5
}
```
#### Response
200 OK

### 5. Удаление книги
`DELETE http://localhost:8080/books/1`

#### Response
200 OK

### 6. Отметить книгу как отсутствующую
`POST http://localhost:8080/books/1/mark-out-of-stock`

#### Response 
200 OK

### 7. Получить рекомендации
`GET http://localhost:8080/books/recommend`


#### Respose
```json
[
    {
        "id": 18,
        "title": "Три товарища",
        "author": "Эрих Мария Ремарк",
        "year": 1936,
        "isbn": "978-5-17-080129-1",
        "outOfStock": false,
        "rating": 10,
        "CreatedAt": "0001-01-01T00:00:00Z",
        "UpdatedAt": "0001-01-01T00:00:00Z"
    },
    {
        "id": 15,
        "title": "Преступление и наказание",
        "author": "Федор Достоевский",
        "year": 1866,
        "isbn": "978-5-17-082645-4",
        "outOfStock": false,
        "rating": 9,
        "CreatedAt": "0001-01-01T00:00:00Z",
        "UpdatedAt": "0001-01-01T00:00:00Z"
    },
    {
        "id": 19,
        "title": "О дивный новый мир",
        "author": "Олдос Хаксли",
        "year": 1932,
        "isbn": "978-5-04-098933-3",
        "outOfStock": false,
        "rating": 7,
        "CreatedAt": "0001-01-01T00:00:00Z",
        "UpdatedAt": "0001-01-01T00:00:00Z"
    },
    {
        "id": 17,
        "title": "1984",
        "author": "Джордж Оруэлл",
        "year": 1949,
        "isbn": "978-5-17-080128-4",
        "outOfStock": false,
        "rating": 6,
        "CreatedAt": "0001-01-01T00:00:00Z",
        "UpdatedAt": "0001-01-01T00:00:00Z"
    },
    {
        "id": 16,
        "title": "Гордость и предубеждение",
        "author": "Джейн Остин",
        "year": 1813,
        "isbn": "978-5-04-098931-9",
        "outOfStock": false,
        "rating": 5,
        "CreatedAt": "0001-01-01T00:00:00Z",
        "UpdatedAt": "0001-01-01T00:00:00Z"
    }
]
```