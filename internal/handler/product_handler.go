package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/product-management-server/internal/model"
	"github.com/product-management-server/internal/service"
)

// ProductHandler handles HTTP requests for product operations.
type ProductHandler struct {
	service service.ProductService
}

// NewProductHandler creates a new ProductHandler with the given service.
func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// CreateProduct handles POST /products.
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req model.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	product, err := h.service.CreateProduct(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProducts handles GET /products with optional ?keyword= query param.
func (h *ProductHandler) GetProducts(c *gin.Context) {
	keyword := c.Query("keyword")

	products, err := h.service.GetProducts(c.Request.Context(), keyword)
	if err != nil {
		h.handleError(c, err)
		return
	}

	if products == nil {
		products = []model.Product{}
	}

	c.JSON(http.StatusOK, products)
}

// GetProductByID handles GET /products/:id.
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := h.parseID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "product not found"})
		return
	}

	product, err := h.service.GetProductByID(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct handles PUT /products/:id.
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := h.parseID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "product not found"})
		return
	}

	var req model.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	product, err := h.service.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct handles DELETE /products/:id.
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := h.parseID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "product not found"})
		return
	}

	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

// parseID extracts and validates the :id path parameter.
// Returns an error if the id is not a positive integer.
func (h *ProductHandler) parseID(c *gin.Context) (int, error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

// handleError maps service errors to appropriate HTTP responses.
func (h *ProductHandler) handleError(c *gin.Context, err error) {
	var validationErr *model.ValidationError
	var notFoundErr *model.NotFoundError

	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, gin.H{"message": validationErr.Message})
		return
	}

	if errors.As(err, &notFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"message": notFoundErr.Message})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
}
