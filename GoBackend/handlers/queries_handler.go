package handlers

import (
	"fmt"
	"frontdesk/models"
	"frontdesk/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QueriesHandler interface {
	CreateQuery(c *gin.Context)
	GetQueries(c *gin.Context)
	ResolveQuery(c *gin.Context)
	GetFAQs(c *gin.Context)
}

type queriesHandler struct {
	queriesService services.QueriesService
}

func NewQueriesHandler(queriesService services.QueriesService) *queriesHandler {
	return &queriesHandler{queriesService: queriesService}
}

func (h *queriesHandler) CreateQuery(c *gin.Context) {
	var req models.Query
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("Invalid request type: %w", err).Error()})
		return
	}

	if !isValidStatus(req.QueryStatus) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query_status"})
		return
	}

	if err := h.queriesService.CreateQuery(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Query created successfully"})
}

func (h *queriesHandler) GetQueries(c *gin.Context) {
	val, err := h.queriesService.GetQueries()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"Success: ": val})
}

func (h *queriesHandler) ResolveQuery(c *gin.Context) {
	// Implementation for resolving a query can be added here
	id := c.Param("id")
	id_number, _ := strconv.Atoi(id)
	var req models.QueryStatus
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("Invalid request type: %w", err).Error()})
		return
	}

	if !isValidStatus(req.QueryStatus) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query_status"})
		return
	}

	err := h.queriesService.ResolveQuery(&req, id_number)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Query response updated successfully"})
}

func (h *queriesHandler) GetFAQs(c *gin.Context) {
	faqs, err := h.queriesService.GetFAQs()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"faqs": faqs})
}

func isValidStatus(s models.Status) bool {
	return s == 0 || // Pending
		s == 1 || // Resolved
		s == 2 // Unresolved
}
