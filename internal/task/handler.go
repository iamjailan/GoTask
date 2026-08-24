package task

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct{ service Service }

type taskRequest struct {
	Title     string `json:"title" binding:"required,min=1,max=255"`
	Completed bool   `json:"completed"`
}

type taskResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewHandler(service Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	routes := router.Group("/api/v1/tasks")
	routes.POST("", h.create)
	routes.GET("", h.list)
	routes.GET("/:id", h.get)
	routes.PUT("/:id", h.update)
	routes.DELETE("/:id", h.delete)
}

func (h *Handler) create(c *gin.Context) {
	var req taskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required and must be 1-255 characters"})
		return
	}
	model, err := h.service.Create(c.Request.Context(), CreateInput{Title: req.Title})
	if err != nil {
		h.serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response(model))
}

func (h *Handler) list(c *gin.Context) {
	models, err := h.service.List(c.Request.Context())
	if err != nil {
		h.serverError(c, err)
		return
	}
	responses := make([]taskResponse, 0, len(models))
	for _, model := range models {
		responses = append(responses, response(model))
	}
	c.JSON(http.StatusOK, responses)
}

func (h *Handler) get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	model, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, response(model))
}

func (h *Handler) update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req taskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required and must be 1-255 characters"})
		return
	}
	model, err := h.service.Update(c.Request.Context(), id, UpdateInput{Title: req.Title, Completed: req.Completed})
	if err != nil {
		h.writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, response(model))
}

func (h *Handler) delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.writeServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return 0, false
	}
	return uint(id), true
}

func (h *Handler) writeServiceError(c *gin.Context, err error) {
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	h.serverError(c, err)
}

func (h *Handler) serverError(c *gin.Context, _ error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}

func response(model Model) taskResponse {
	return taskResponse{
		ID:        model.ID,
		Title:     model.Title,
		Completed: model.Completed,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
