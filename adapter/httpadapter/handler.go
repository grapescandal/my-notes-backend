package httpadapter

import (
	"net/http"

	"my-note/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, svc *usecase.NoteService) {
	r.Use(cors.Default())

	api := r.Group("/api")

	api.POST("/notes", func(c *gin.Context) {
		var req struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		n, err := svc.CreateNote(req.Title, req.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, n)
	})

	api.GET("/notes", func(c *gin.Context) {
		notes, err := svc.ListNotes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, notes)
	})

	api.GET("/notes/:id", func(c *gin.Context) {
		id := c.Param("id")
		n, err := svc.GetNote(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
			return
		}
		c.JSON(http.StatusOK, n)
	})
}
