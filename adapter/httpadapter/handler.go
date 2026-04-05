package httpadapter

import (
	"net/http"
	"time"

	"my-note/domain"
	"my-note/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, svc *usecase.NoteService) {
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(config))

	api := r.Group("/api")

	api.POST("/notes", func(c *gin.Context) {
		var req struct {
			ID      *string `json:"id" binding:"omitempty"`
			Title   string  `json:"title" binding:"required"`
			Content *string `json:"content" binding:"omitempty"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		note := domain.Note{
			Title: req.Title,
		}
		if req.Content != nil {
			note.Content = *req.Content
		}
		if req.ID != nil {
			note.ID = *req.ID
		}

		n, created, err := svc.SaveNote(note)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if created {
			c.JSON(http.StatusCreated, n)
			return
		}
		c.JSON(http.StatusOK, n)
	})

	api.GET("/notes", func(c *gin.Context) {
		notes, err := svc.ListNotes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// return only id and title, already sorted by service
		out := make([]struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		}, 0, len(notes))
		for _, n := range notes {
			out = append(out, struct {
				ID    string `json:"id"`
				Title string `json:"title"`
			}{ID: n.ID, Title: n.Title})
		}
		c.JSON(http.StatusOK, out)
	})

	api.GET("/notes/:id", func(c *gin.Context) {
		var uri struct {
			ID string `uri:"id" binding:"required"`
		}
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}
		n, err := svc.GetNote(uri.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
			return
		}
		c.JSON(http.StatusOK, n)
	})

	api.DELETE("/notes/:id", func(c *gin.Context) {
		var uri struct {
			ID string `uri:"id" binding:"required"`
		}
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}
		err := svc.DeleteNote(uri.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
			return
		}
		c.Status(http.StatusNoContent)
	})
}
