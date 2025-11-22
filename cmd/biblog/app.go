package main

import (
	"fmt"
	"os"

	"bibliography_log/internal/infrastructure"
	"bibliography_log/internal/service"
)

// App holds the application dependencies.
type App struct {
	BibService    *service.BibliographyService
	ReviewService *service.ReviewService
}

// NewApp initializes the application and its dependencies.
func NewApp() (*App, error) {
	// Ensure data directory exists
	dataDir := "data"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.Mkdir(dataDir, 0755); err != nil {
			return nil, fmt.Errorf("error creating data directory: %w", err)
		}
	}

	// Initialize Repositories
	bibRepo := infrastructure.NewCSVBibliographyRepository(dataDir + "/bibliographies.csv")
	classRepo := infrastructure.NewCSVBibClassificationRepository(dataDir + "/classifications.csv")
	reviewRepo := infrastructure.NewCSVReviewRepository(dataDir + "/reviews.csv")

	// Initialize Service
	bibSvc := service.NewBibliographyService(bibRepo, classRepo)
	reviewSvc := service.NewReviewService(reviewRepo, bibRepo)

	return &App{
		BibService:    bibSvc,
		ReviewService: reviewSvc,
	}, nil
}
