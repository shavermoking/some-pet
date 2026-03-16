package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"some-pet/internal/models"
	"some-pet/internal/service"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockBooksRepository struct {
	mock.Mock
}

func (m *MockBooksRepository) Create(ctx context.Context, book models.Book) (*models.Book, error) {
	args := m.Called(ctx, book)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBooksRepository) GetAll(ctx context.Context) ([]models.Book, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Book), args.Error(1)
}

func (m *MockBooksRepository) GetByID(ctx context.Context, id int) (models.Book, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Book), args.Error(1)
}

func (m *MockBooksRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBooksRepository) Update(ctx context.Context, id int, input models.UpdateBook) error {
	args := m.Called(ctx, id, input)
	return args.Error(0)
}

func (m *MockBooksRepository) MarkOutOfStock(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBooksRepository) GetRecommend(ctx context.Context) ([]models.Book, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Book), args.Error(1)
}

func setupRouter(handler *Books) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	books := r.Group("/api/books")
	{
		books.POST("/", handler.Create)
		books.GET("/", handler.GetAll)
		books.GET("/:id", handler.GetByID)
		books.DELETE("/:id", handler.Delete)
		books.PUT("/:id", handler.Update)
		books.PATCH("/:id/out-of-stock", handler.MarkOutOfStock)
		books.GET("/recommend", handler.GetRecommend)
	}

	return r
}

func TestBooks_Create(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(mockRepo *MockBooksRepository)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "success",
			requestBody: models.Book{
				Title:  "Test Book",
				Author: "Test Author",
				Year:   2024,
				ISBN:   "123-456789",
				Rating: 5,
			},
			mockSetup: func(mockRepo *MockBooksRepository) {
				mockRepo.On("Create", mock.Anything, models.Book{
					Title:  "Test Book",
					Author: "Test Author",
					Year:   2024,
					ISBN:   "123-456789",
					Rating: 5,
				}).Return(&models.Book{
					ID:     1,
					Title:  "Test Book",
					Author: "Test Author",
					Year:   2024,
					ISBN:   "123-456789",
					Rating: 5,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name: "invalid json - wrong types",
			requestBody: map[string]interface{}{
				"title": 123,
				"year":  "invalid",
			},
			mockSetup:      func(mockRepo *MockBooksRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "некорректный формат данных",
		},
		{
			name: "service error - duplicate ISBN",
			requestBody: models.Book{
				Title:  "Test Book",
				Author: "Test Author",
				Year:   2024,
				ISBN:   "123-456789",
				Rating: 5,
			},
			mockSetup: func(mockRepo *MockBooksRepository) {
				mockRepo.On("Create", mock.Anything, models.Book{
					Title:  "Test Book",
					Author: "Test Author",
					Year:   2024,
					ISBN:   "123-456789",
					Rating: 5,
				}).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "книга с таким ISBN уже существует или данные некорректны",
		},
		{
			name: "validation error - missing required fields",
			requestBody: models.Book{
				Title: "Test Book",
			},
			mockSetup:      func(mockRepo *MockBooksRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "некорректный формат данных",
		},
		{
			name:           "invalid json - malformed",
			requestBody:    "not a json",
			mockSetup:      func(mockRepo *MockBooksRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "некорректный формат данных",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBooksRepository)
			tt.mockSetup(mockRepo)

			serviceInstance := service.NewBooks(mockRepo)
			handler := NewBooks(serviceInstance)
			router := setupRouter(handler)

			var jsonBody []byte
			var err error

			if strBody, ok := tt.requestBody.(string); ok && strBody == "not a json" {
				jsonBody = []byte("{invalid json")
			} else {
				jsonBody, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			}

			req, err := http.NewRequest(http.MethodPost, "/api/books/", bytes.NewBuffer(jsonBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			if tt.expectedError != "" {
				errorMsg, exists := response["error"]
				assert.True(t, exists, "Response should contain error field")
				assert.Contains(t, errorMsg, tt.expectedError,
					"Error message should contain expected substring")
			} else {
				assert.Contains(t, response, "id")
				assert.NotEmpty(t, response["id"])
			}

			if tt.name == "success" || tt.name == "service error - duplicate ISBN" {
				mockRepo.AssertExpectations(t)
			} else {
				mockRepo.AssertNotCalled(t, "Create")
			}
		})
	}
}

func TestBooks_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		mockBooks      []models.Book
		mockError      error
		expectedStatus int
	}{
		{
			name: "success",
			mockBooks: []models.Book{
				{ID: 1, Title: "Book 1", Author: "Author 1", Year: 2024, ISBN: "123", Rating: 5},
				{ID: 2, Title: "Book 2", Author: "Author 2", Year: 2023, ISBN: "456", Rating: 4},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty list",
			mockBooks:      []models.Book{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "service error",
			mockBooks:      []models.Book{},
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBooksRepository)
			serviceInstance := service.NewBooks(mockRepo)
			handler := NewBooks(serviceInstance)
			router := setupRouter(handler)

			mockRepo.On("GetAll", mock.Anything).Return(tt.mockBooks, tt.mockError)

			req, _ := http.NewRequest(http.MethodGet, "/api/books/", nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.mockError == nil {
				var response []models.Book
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockBooks, response)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBooks_GetByID(t *testing.T) {
	expectedBook := models.Book{
		ID:     1,
		Title:  "Test Book",
		Author: "Test Author",
		Year:   2024,
		ISBN:   "123-456789",
		Rating: 5,
	}

	tests := []struct {
		name           string
		id             string
		mockBook       models.Book
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "success",
			id:             "1",
			mockBook:       expectedBook,
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid id",
			id:             "abc",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid id",
			},
		},
		{
			name:           "book not found",
			id:             "999",
			mockBook:       models.Book{},
			mockError:      models.ErrBookNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "book not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBooksRepository)
			serviceInstance := service.NewBooks(mockRepo)
			handler := NewBooks(serviceInstance)
			router := setupRouter(handler)

			if tt.name != "invalid id" {
				mockRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int")).Return(tt.mockBook, tt.mockError)
			}

			req, _ := http.NewRequest(http.MethodGet, "/api/books/"+tt.id, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}

			if tt.name != "invalid id" {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestBooks_Delete(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "success",
			id:             "1",
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "invalid id",
			id:             "abc",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid id",
			},
		},
		{
			name:           "service error",
			id:             "1",
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "database error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBooksRepository)
			serviceInstance := service.NewBooks(mockRepo)
			handler := NewBooks(serviceInstance)
			router := setupRouter(handler)

			if tt.name != "invalid id" {
				mockRepo.On("Delete", mock.Anything, mock.AnythingOfType("int")).Return(tt.mockError)
			}

			req, _ := http.NewRequest(http.MethodDelete, "/api/books/"+tt.id, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}

			if tt.name != "invalid id" {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestBooks_Update(t *testing.T) {
	title := "Updated Title"
	author := "Updated Author"
	year := 2025
	isbn := "987-654321"
	outOfStock := true
	rating := 4

	tests := []struct {
		name           string
		id             string
		requestBody    interface{}
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success",
			id:   "1",
			requestBody: models.UpdateBook{
				Title:      &title,
				Author:     &author,
				Year:       &year,
				ISBN:       &isbn,
				OutOfStock: &outOfStock,
				Rating:     &rating,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "success with partial update",
			id:   "1",
			requestBody: models.UpdateBook{
				Title:  &title,
				Rating: &rating,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid id",
			id:   "abc",
			requestBody: models.UpdateBook{
				Title: &title,
			},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid id",
			},
		},
		{
			name: "invalid json",
			id:   "1",
			requestBody: map[string]interface{}{
				"title": 123,
				"year":  "invalid",
			},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			id:   "1",
			requestBody: models.UpdateBook{
				Title: &title,
			},
			mockError:      errors.New("update failed"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "update failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBooksRepository)
			serviceInstance := service.NewBooks(mockRepo)
			handler := NewBooks(serviceInstance)
			router := setupRouter(handler)

			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPut, "/api/books/"+tt.id, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.name == "success" || tt.name == "success with partial update" {
				if updateBook, ok := tt.requestBody.(models.UpdateBook); ok {
					mockRepo.On("Update", mock.Anything, 1, updateBook).Return(tt.mockError)
				}
			} else if tt.name == "service error" {
				if updateBook, ok := tt.requestBody.(models.UpdateBook); ok {
					mockRepo.On("Update", mock.Anything, 1, updateBook).Return(tt.mockError)
				}
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}

			if tt.name == "success" || tt.name == "success with partial update" || tt.name == "service error" {
				if _, ok := tt.requestBody.(models.UpdateBook); ok {
					mockRepo.AssertExpectations(t)
				}
			}
		})
	}
}

func TestBooks_MarkOutOfStock(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "success",
			id:             "1",
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid id",
			id:             "abc",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid id",
			},
		},
		{
			name:           "book not found",
			id:             "999",
			mockError:      models.ErrBookNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "book not found",
			},
		},
		{
			name:           "service error",
			id:             "1",
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "database error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBooksRepository)
			serviceInstance := service.NewBooks(mockRepo)
			handler := NewBooks(serviceInstance)
			router := setupRouter(handler)

			if tt.name != "invalid id" {
				mockRepo.On("MarkOutOfStock", mock.Anything, mock.AnythingOfType("int")).Return(tt.mockError)
			}

			req, _ := http.NewRequest(http.MethodPatch, "/api/books/"+tt.id+"/out-of-stock", nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}

			if tt.name != "invalid id" {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestBooks_GetRecommend(t *testing.T) {
	tests := []struct {
		name           string
		mockBooks      []models.Book
		mockError      error
		expectedStatus int
	}{
		{
			name: "success",
			mockBooks: []models.Book{
				{ID: 1, Title: "Recommended 1", Author: "Author 1", Year: 2024, ISBN: "123", Rating: 5},
				{ID: 2, Title: "Recommended 2", Author: "Author 2", Year: 2023, ISBN: "456", Rating: 5},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty list",
			mockBooks:      []models.Book{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "service error",
			mockBooks:      []models.Book{},
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBooksRepository)
			serviceInstance := service.NewBooks(mockRepo)
			handler := NewBooks(serviceInstance)
			router := setupRouter(handler)

			mockRepo.On("GetRecommend", mock.Anything).Return(tt.mockBooks, tt.mockError)

			req, _ := http.NewRequest(http.MethodGet, "/api/books/recommend", nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.mockError == nil {
				var response []models.Book
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockBooks, response)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
