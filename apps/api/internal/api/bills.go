// Package api exposes the HTTP CRUD API for the 415 service.
package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ofgrenudo/415/apps/api/internal/models"
	"github.com/ofgrenudo/415/apps/api/internal/store"
)

// BillsHandler serves CRUD requests for bills.
type BillsHandler struct {
	store *store.BillStore
}

// NewBillsHandler returns a BillsHandler backed by s.
func NewBillsHandler(s *store.BillStore) *BillsHandler {
	return &BillsHandler{store: s}
}

// Register wires the bills routes onto router.
func (h *BillsHandler) Register(router gin.IRouter) {
	router.GET("/bills", h.list)
	router.POST("/bills", h.create)
	router.GET("/bills/:id", h.get)
	router.PUT("/bills/:id", h.update)
	router.DELETE("/bills/:id", h.delete)
}

func (h *BillsHandler) list(c *gin.Context) {
	bills, err := h.store.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorBody(err))
		return
	}
	c.JSON(http.StatusOK, bills)
}

func (h *BillsHandler) get(c *gin.Context) {
	id, err := idFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorBody(err))
		return
	}

	bill, err := h.store.Get(c.Request.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		c.JSON(http.StatusNotFound, errorBody(err))
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorBody(err))
		return
	}
	c.JSON(http.StatusOK, bill)
}

func (h *BillsHandler) create(c *gin.Context) {
	var in models.Bill
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, errorBody(err))
		return
	}
	if in.Name == "" {
		c.JSON(http.StatusBadRequest, errorBody(errors.New("name is required")))
		return
	}

	bill, err := h.store.Create(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorBody(err))
		return
	}
	c.JSON(http.StatusCreated, bill)
}

func (h *BillsHandler) update(c *gin.Context) {
	id, err := idFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorBody(err))
		return
	}

	var in models.Bill
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, errorBody(err))
		return
	}
	if in.Name == "" {
		c.JSON(http.StatusBadRequest, errorBody(errors.New("name is required")))
		return
	}

	bill, err := h.store.Update(c.Request.Context(), id, in)
	if errors.Is(err, store.ErrNotFound) {
		c.JSON(http.StatusNotFound, errorBody(err))
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorBody(err))
		return
	}
	c.JSON(http.StatusOK, bill)
}

func (h *BillsHandler) delete(c *gin.Context) {
	id, err := idFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorBody(err))
		return
	}

	if err := h.store.Delete(c.Request.Context(), id); errors.Is(err, store.ErrNotFound) {
		c.JSON(http.StatusNotFound, errorBody(err))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, errorBody(err))
		return
	}
	c.Status(http.StatusNoContent)
}

func idFromParam(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
}

func errorBody(err error) gin.H {
	return gin.H{"error": err.Error()}
}
