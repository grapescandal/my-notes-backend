package main

import (
	"log"
	"os"

	"my-note/adapter/httpadapter"
	"my-note/adapter/memory"
	"my-note/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	repo := memory.NewMemoryRepo()
	svc := usecase.NewNoteService(repo)

	r := gin.Default()
	httpadapter.RegisterRoutes(r, svc)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
