package handler

import (
	"calendar/internal/models"
	"calendar/internal/service"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger is a middleware that logs each request
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		method := c.Request.Method
		path := c.Request.URL.Path

		c.Next()

		status := c.Writer.Status()
		fmt.Printf("[%s] %s %s -> %d (%v)\n",
			start.Format("2006-01-02 15:04:05"),
			method, path, status, time.Since(start))
	}
}

// Handler is a struct that contains a service
type Handler struct {
	service service.Service
}

// NewHandler creates a new handler
func NewHandler(service service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) createEvent(c *gin.Context) {
	e := new(models.Event)

	if err := c.ShouldBindJSON(e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := e.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event " + err.Error(),
		})
		return
	}

	h.service.CreateEvent(e)

	c.JSON(http.StatusOK, gin.H{
		"result": e,
	})
}

func (h *Handler) updateEvent(c *gin.Context) {
	e := new(models.Event)

	if err := c.ShouldBindJSON(e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.UpdateEvent(e)
	if errors.Is(err, models.ErrNotFound) {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": err.Error(),
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": e,
	})
}

func (h *Handler) deleteEvent(c *gin.Context) {
	e := new(models.Event)

	if err := c.ShouldBindJSON(e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.DeleteEvent(e.ID)
	if errors.Is(err, models.ErrNotFound) {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": err.Error(),
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": e,
	})
}

func (h *Handler) getEventsSince(c *gin.Context, since models.Date) {
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user_id" + err.Error(),
		})
		return
	}

	date, err := time.Parse("2006-01-02", c.Query("date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date" + err.Error(),
		})
		return
	}

	events := h.service.GetSince(userID, date, since)

	c.JSON(http.StatusOK, gin.H{
		"result": events.Events,
	})
}

// RegisterRoutes registers the routes for the handler
//
//	Post /create_event
//	Post /update_event
//	Post /delete_event
//	Get /events_for_day
//	Get /events_for_week
//	Get /events_for_month
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/create_event", h.createEvent)
	r.POST("/update_event", h.updateEvent)
	r.POST("/delete_event", h.deleteEvent)

	r.GET("/events_for_day", func(c *gin.Context) {
		h.getEventsSince(c, models.Day)
	})
	r.GET("/events_for_week", func(c *gin.Context) {
		h.getEventsSince(c, models.Week)
	})
	r.GET("/events_for_month", func(c *gin.Context) {
		h.getEventsSince(c, models.Month)
	})
}
