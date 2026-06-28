package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/product-management-server/internal/model"
	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"
)

// --- Mock Service ---

type mockProductService struct {
	createProductFn  func(ctx context.Context, req model.CreateProductRequest) (*model.Product, error)
	getProductsFn    func(ctx context.Context, keyword string) ([]model.Product, error)
	getProductByIDFn func(ctx context.Context, id int) (*model.Product, error)
	updateProductFn  func(ctx context.Context, id int, req model.UpdateProductRequest) (*model.Product, error)
	deleteProductFn  func(ctx context.Context, id int) error
}

func (m *mockProductService) CreateProduct(ctx context.Context, req model.CreateProductRequest) (*model.Product, error) {
	if m.createProductFn != nil {
		return m.createProductFn(ctx, req)
	}
	return nil, nil
}

func (m *mockProductService) GetProducts(ctx context.Context, keyword string) ([]model.Product, error) {
	if m.getProductsFn != nil {
		return m.getProductsFn(ctx, keyword)
	}
	return nil, nil
}

func (m *mockProductService) GetProductByID(ctx context.Context, id int) (*model.Product, error) {
	if m.getProductByIDFn != nil {
		return m.getProductByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockProductService) UpdateProduct(ctx context.Context, id int, req model.UpdateProductRequest) (*model.Product, error) {
	if m.updateProductFn != nil {
		return m.updateProductFn(ctx, id, req)
	}
	return nil, nil
}

func (m *mockProductService) DeleteProduct(ctx context.Context, id int) error {
	if m.deleteProductFn != nil {
		return m.deleteProductFn(ctx, id)
	}
	return nil
}

// --- Helper ---

// setupRouter registers all product routes directly to avoid import cycles
// between the handler and routes packages.
func setupRouter(svc *mockProductService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := NewProductHandler(svc)
	router.POST("/products", h.CreateProduct)
	router.GET("/products", h.GetProducts)
	router.GET("/products/:id", h.GetProductByID)
	router.PUT("/products/:id", h.UpdateProduct)
	router.DELETE("/products/:id", h.DeleteProduct)
	return router
}

// --- Unit Tests ---

func TestCreateProduct_MalformedJSON(t *testing.T) {
	svc := &mockProductService{}
	router := setupRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "invalid request body", body["message"])
}

func TestCreateProduct_ValidationError(t *testing.T) {
	svc := &mockProductService{
		createProductFn: func(ctx context.Context, req model.CreateProductRequest) (*model.Product, error) {
			return nil, &model.ValidationError{Message: "name is required"}
		},
	}
	router := setupRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"name":"","price":10.0,"quantity":5}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "name is required", body["message"])
}

func TestGetProductByID_NonIntegerID(t *testing.T) {
	svc := &mockProductService{}
	router := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/products/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "product not found", body["message"])
}

func TestGetProductByID_NegativeID(t *testing.T) {
	svc := &mockProductService{}
	router := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/products/-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "product not found", body["message"])
}

func TestGetProductByID_NotFoundError(t *testing.T) {
	svc := &mockProductService{
		getProductByIDFn: func(ctx context.Context, id int) (*model.Product, error) {
			return nil, &model.NotFoundError{Message: "product not found"}
		},
	}
	router := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/products/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "product not found", body["message"])
}

func TestHandler_GenericError_Returns500(t *testing.T) {
	svc := &mockProductService{
		getProductByIDFn: func(ctx context.Context, id int) (*model.Product, error) {
			return nil, errors.New("database connection failed")
		},
	}
	router := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "internal server error", body["message"])
}

func TestHandler_ContentTypeJSON_AllResponses(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		body   string
		svc    *mockProductService
	}{
		{
			name:   "CreateProduct success",
			method: http.MethodPost,
			path:   "/products",
			body:   `{"name":"Test Product","price":10.0,"quantity":5}`,
			svc: &mockProductService{
				createProductFn: func(ctx context.Context, req model.CreateProductRequest) (*model.Product, error) {
					price := 10.0
					qty := 5
					return &model.Product{ID: 1, Name: req.Name, Price: price, Quantity: qty}, nil
				},
			},
		},
		{
			name:   "CreateProduct error",
			method: http.MethodPost,
			path:   "/products",
			body:   "{bad json",
			svc:    &mockProductService{},
		},
		{
			name:   "GetProducts success",
			method: http.MethodGet,
			path:   "/products",
			body:   "",
			svc: &mockProductService{
				getProductsFn: func(ctx context.Context, keyword string) ([]model.Product, error) {
					return []model.Product{}, nil
				},
			},
		},
		{
			name:   "GetProductByID not found",
			method: http.MethodGet,
			path:   "/products/abc",
			body:   "",
			svc:    &mockProductService{},
		},
		{
			name:   "DeleteProduct success",
			method: http.MethodDelete,
			path:   "/products/1",
			body:   "",
			svc: &mockProductService{
				deleteProductFn: func(ctx context.Context, id int) error {
					return nil
				},
			},
		},
		{
			name:   "Internal server error",
			method: http.MethodGet,
			path:   "/products/1",
			body:   "",
			svc: &mockProductService{
				getProductByIDFn: func(ctx context.Context, id int) (*model.Product, error) {
					return nil, errors.New("unexpected error")
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router := setupRouter(tc.svc)

			var bodyReader *strings.Reader
			if tc.body != "" {
				bodyReader = strings.NewReader(tc.body)
			} else {
				bodyReader = strings.NewReader("")
			}

			req := httptest.NewRequest(tc.method, tc.path, bodyReader)
			if tc.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Contains(t, w.Header().Get("Content-Type"), "application/json",
				"Content-Type should be application/json for %s", tc.name)
		})
	}
}

func TestGetProducts_EmptyArray(t *testing.T) {
	svc := &mockProductService{
		getProductsFn: func(ctx context.Context, keyword string) ([]model.Product, error) {
			return nil, nil
		},
	}
	router := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	assert.Equal(t, "[]", strings.TrimSpace(w.Body.String()))
}

func TestDeleteProduct_Success(t *testing.T) {
	svc := &mockProductService{
		deleteProductFn: func(ctx context.Context, id int) error {
			return nil
		},
	}
	router := setupRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "product deleted successfully", body["message"])
}

func TestRouteRegistration(t *testing.T) {
	svc := &mockProductService{
		createProductFn: func(ctx context.Context, req model.CreateProductRequest) (*model.Product, error) {
			return &model.Product{ID: 1, Name: req.Name}, nil
		},
		getProductsFn: func(ctx context.Context, keyword string) ([]model.Product, error) {
			return []model.Product{}, nil
		},
		getProductByIDFn: func(ctx context.Context, id int) (*model.Product, error) {
			return &model.Product{ID: id, Name: "Test"}, nil
		},
		updateProductFn: func(ctx context.Context, id int, req model.UpdateProductRequest) (*model.Product, error) {
			return &model.Product{ID: id, Name: req.Name}, nil
		},
		deleteProductFn: func(ctx context.Context, id int) error {
			return nil
		},
	}
	router := setupRouter(svc)

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
	}{
		{"POST /products", http.MethodPost, "/products", `{"name":"Test","price":10.0,"quantity":5}`, http.StatusCreated},
		{"GET /products", http.MethodGet, "/products", "", http.StatusOK},
		{"GET /products/:id", http.MethodGet, "/products/1", "", http.StatusOK},
		{"PUT /products/:id", http.MethodPut, "/products/1", `{"name":"Updated","price":20.0,"quantity":10}`, http.StatusOK},
		{"DELETE /products/:id", http.MethodDelete, "/products/1", "", http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var bodyReader *strings.Reader
			if tc.body != "" {
				bodyReader = strings.NewReader(tc.body)
			} else {
				bodyReader = strings.NewReader("")
			}

			req := httptest.NewRequest(tc.method, tc.path, bodyReader)
			if tc.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Should not return 404 (route not found) - confirms route is registered
			assert.Equal(t, tc.expectedStatus, w.Code, "Route %s %s should be registered", tc.method, tc.path)
		})
	}
}


// --- Property 8: Error responses have consistent structure ---

// errorResponseBody represents the expected JSON error response structure.
type errorResponseBody struct {
	Message string `json:"message"`
}

// TestProperty8_ErrorResponsesHaveConsistentStructure tests Property 8:
// Error responses have consistent structure.
// For any request that triggers an error (validation error, not found, or internal error),
// the JSON response body SHALL have the structure {"message": "<description>"} with a
// non-empty message string.
func TestProperty8_ErrorResponsesHaveConsistentStructure(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Choose an error scenario to test
		scenario := rapid.IntRange(0, 4).Draw(t, "scenario")

		var (
			recorder       *httptest.ResponseRecorder
			expectedStatus int
		)

		switch scenario {
		case 0:
			// Scenario 1: Malformed JSON body → 400 with "invalid request body"
			recorder, expectedStatus = prop8_scenarioMalformedJSON(t)
		case 1:
			// Scenario 2: ValidationError from service → 400 with error message
			recorder, expectedStatus = prop8_scenarioValidationError(t)
		case 2:
			// Scenario 3: NotFoundError from service → 404 with error message
			recorder, expectedStatus = prop8_scenarioNotFoundError(t)
		case 3:
			// Scenario 4: Generic errors from service → 500 with "internal server error"
			recorder, expectedStatus = prop8_scenarioGenericError(t)
		case 4:
			// Scenario 5: Invalid :id path params → 404 with "product not found"
			recorder, expectedStatus = prop8_scenarioInvalidID(t)
		}

		// Assert: status code matches expected
		if recorder.Code != expectedStatus {
			t.Fatalf("expected HTTP status %d, got %d (scenario %d)", expectedStatus, recorder.Code, scenario)
		}

		// Assert: response body is valid JSON with structure {"message": "<non-empty string>"}
		var resp errorResponseBody
		err := json.Unmarshal(recorder.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("response body is not valid JSON: %s (body: %s)", err, recorder.Body.String())
		}

		// Assert: message field is non-empty
		if resp.Message == "" {
			t.Fatalf("error message should be non-empty, got empty string")
		}

		// Assert: response body has ONLY the "message" field (consistent structure)
		var rawMap map[string]interface{}
		err = json.Unmarshal(recorder.Body.Bytes(), &rawMap)
		if err != nil {
			t.Fatalf("response body should unmarshal to map: %s", err)
		}
		if len(rawMap) != 1 {
			t.Fatalf("error response should have exactly 1 field, got %d: %v", len(rawMap), rawMap)
		}
		msgVal, hasMessage := rawMap["message"]
		if !hasMessage {
			t.Fatalf("error response must have 'message' field, got: %v", rawMap)
		}
		if _, isString := msgVal.(string); !isString {
			t.Fatalf("'message' field must be a string, got %T", msgVal)
		}
	})
}

// prop8_scenarioMalformedJSON sends a POST /products request with malformed JSON body.
func prop8_scenarioMalformedJSON(t *rapid.T) (*httptest.ResponseRecorder, int) {
	svc := &mockProductService{}
	router := setupRouter(svc)

	// Generate malformed JSON strings that are definitely not valid JSON
	malformedType := rapid.IntRange(0, 3).Draw(t, "malformedType")
	var body string
	switch malformedType {
	case 0:
		// Missing closing brace
		body = `{"name": "test"`
	case 1:
		// Random garbage text
		body = rapid.StringMatching(`[a-zA-Z0-9!@#$%]{1,50}`).Draw(t, "garbageBody")
	case 2:
		// Broken key-value pairs
		body = `{name: invalid}`
	case 3:
		// Truncated JSON
		body = `{"name": "product", "price":`
	}

	req, _ := http.NewRequest("POST", "/products", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	return recorder, http.StatusBadRequest
}

// prop8_scenarioValidationError configures the mock service to return a ValidationError.
func prop8_scenarioValidationError(t *rapid.T) (*httptest.ResponseRecorder, int) {
	// Generate a non-empty validation message
	msg := rapid.StringMatching(`[a-z ]{3,50}`).Draw(t, "validationMsg")

	svc := &mockProductService{
		createProductFn: func(ctx context.Context, req model.CreateProductRequest) (*model.Product, error) {
			return nil, &model.ValidationError{Message: msg}
		},
	}
	router := setupRouter(svc)

	// Send a valid JSON body (the service mock will return the validation error)
	body := `{"name":"ValidName","price":10.0,"quantity":5}`
	req, _ := http.NewRequest("POST", "/products", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	return recorder, http.StatusBadRequest
}

// prop8_scenarioNotFoundError configures the mock service to return a NotFoundError.
func prop8_scenarioNotFoundError(t *rapid.T) (*httptest.ResponseRecorder, int) {
	// Generate a non-empty not-found message
	msg := rapid.StringMatching(`[a-z ]{3,50}`).Draw(t, "notFoundMsg")

	svc := &mockProductService{
		getProductByIDFn: func(ctx context.Context, id int) (*model.Product, error) {
			return nil, &model.NotFoundError{Message: msg}
		},
	}
	router := setupRouter(svc)

	// Use a valid ID so parseID passes, then the service returns NotFoundError
	id := rapid.IntRange(1, 999999).Draw(t, "validID")
	url := fmt.Sprintf("/products/%d", id)

	req, _ := http.NewRequest("GET", url, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	return recorder, http.StatusNotFound
}

// prop8_scenarioGenericError configures the mock service to return a generic (non-typed) error.
func prop8_scenarioGenericError(t *rapid.T) (*httptest.ResponseRecorder, int) {
	// Generate a random error message (this won't be exposed to client)
	errMsg := rapid.StringMatching(`[a-z ]{3,50}`).Draw(t, "genericErrMsg")

	// Choose which endpoint to test generic errors on
	endpoint := rapid.IntRange(0, 2).Draw(t, "endpoint")

	var svc *mockProductService
	var req *http.Request

	switch endpoint {
	case 0:
		// CreateProduct with generic error
		svc = &mockProductService{
			createProductFn: func(ctx context.Context, r model.CreateProductRequest) (*model.Product, error) {
				return nil, errors.New(errMsg)
			},
		}
		body := `{"name":"ValidName","price":10.0,"quantity":5}`
		req, _ = http.NewRequest("POST", "/products", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
	case 1:
		// GetProductByID with generic error
		svc = &mockProductService{
			getProductByIDFn: func(ctx context.Context, id int) (*model.Product, error) {
				return nil, errors.New(errMsg)
			},
		}
		id := rapid.IntRange(1, 999999).Draw(t, "id")
		req, _ = http.NewRequest("GET", fmt.Sprintf("/products/%d", id), nil)
	case 2:
		// DeleteProduct with generic error
		svc = &mockProductService{
			deleteProductFn: func(ctx context.Context, id int) error {
				return errors.New(errMsg)
			},
		}
		id := rapid.IntRange(1, 999999).Draw(t, "deleteID")
		req, _ = http.NewRequest("DELETE", fmt.Sprintf("/products/%d", id), nil)
	}

	router := setupRouter(svc)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	return recorder, http.StatusInternalServerError
}

// prop8_scenarioInvalidID sends a request with an invalid :id path parameter.
func prop8_scenarioInvalidID(t *rapid.T) (*httptest.ResponseRecorder, int) {
	svc := &mockProductService{}
	router := setupRouter(svc)

	// Generate invalid IDs: non-numeric strings, zero, or negative numbers
	invalidType := rapid.IntRange(0, 2).Draw(t, "invalidIDType")
	var url string

	switch invalidType {
	case 0:
		// Non-numeric string
		invalidStr := rapid.StringMatching(`[a-zA-Z]{1,10}`).Draw(t, "nonNumericID")
		url = fmt.Sprintf("/products/%s", invalidStr)
	case 1:
		// Zero
		url = "/products/0"
	case 2:
		// Negative number
		negID := rapid.IntRange(-99999, -1).Draw(t, "negativeID")
		url = fmt.Sprintf("/products/%d", negID)
	}

	// Choose which method to use (GET, PUT, DELETE all parse :id)
	method := rapid.SampledFrom([]string{"GET", "PUT", "DELETE"}).Draw(t, "method")

	if method == "PUT" {
		// PUT needs a body but since ID parsing happens first, it doesn't matter
		req, _ := http.NewRequest(method, url, bytes.NewBufferString(`{"name":"test","price":10,"quantity":5}`))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		return recorder, http.StatusNotFound
	}

	req, _ := http.NewRequest(method, url, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	return recorder, http.StatusNotFound
}
