package main

import (
	"bibliography_log/internal/domain"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	app, err := NewApp()
	if err != nil {
		fmt.Printf("Error initializing application: %v\n", err)
		os.Exit(1)
	}

	// Subcommands
	addClassCmd := flag.NewFlagSet("add-class", flag.ExitOnError)
	addBibCmd := flag.NewFlagSet("add-bib", flag.ExitOnError)
	addReviewCmd := flag.NewFlagSet("add-review", flag.ExitOnError)
	updateReviewCmd := flag.NewFlagSet("update-review", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	// Add Class Flags
	addClassCode := addClassCmd.Int("code", 0, "Classification Code Number (e.g. 56)")
	addClassName := addClassCmd.String("name", "", "Classification Name (e.g. Technology)")

	// Add Bib Flags
	addBibTitle := addBibCmd.String("title", "", "Title of the bibliography")
	addBibAuthor := addBibCmd.String("author", "", "Author of the bibliography")
	addBibPublisher := addBibCmd.String("publisher", "", "Publisher of the bibliography")
	addBibType := addBibCmd.String("type", "", "Type (Book, Essay, Video, etc.)")
	addBibClass := addBibCmd.Int("class", 0, "Classification Code Number")
	addBibYear := addBibCmd.Int("year", 0, "Published Year (e.g. 2024)")
	addBibISBN := addBibCmd.String("isbn", "", "ISBN")
	addBibTitleEn := addBibCmd.String("title-en", "", "English translation of title (required if title contains Japanese)")
	addBibAuthorEn := addBibCmd.String("author-en", "", "English translation of author (required if author contains Japanese)")
	addBibIndex := addBibCmd.String("bib-index", "", "Manual BibIndex (overrides auto-generation and bypasses English translation requirements)")

	// Add Review Flags
	addReviewBibIndex := addReviewCmd.String("bib-index", "", "BibIndex of the bibliography to review")
	addReviewGoals := addReviewCmd.String("goals", "", "Goals for reading (required)")
	addReviewSummary := addReviewCmd.String("summary", "", "Summary of the review")

	// Update Review Flags
	updateReviewID := updateReviewCmd.String("review-id", "", "UUID of the review to update (required)")
	updateReviewGoals := updateReviewCmd.String("goals", "", "New goals for reading (optional)")
	updateReviewSummary := updateReviewCmd.String("summary", "", "New summary of the review (optional)")

	// List Flags
	listLimit := listCmd.Int("limit", 100, "Maximum number of items to display (default: 100, 0 for all)")
	listOffset := listCmd.Int("offset", 0, "Number of items to skip (default: 0)")

	if len(os.Args) < 2 {
		fmt.Println("expected 'add-class', 'add-bib', 'add-review', 'update-review' or 'list' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add-class":
		_ = addClassCmd.Parse(os.Args[2:])
		if *addClassCode == 0 || *addClassName == "" {
			fmt.Println("Please provide -code and -name")
			addClassCmd.PrintDefaults()
			os.Exit(1)
		}
		class, err := app.BibService.AddClassification(*addClassCode, *addClassName)
		if err != nil {
			fmt.Printf("Error adding classification: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Classification added: %v\n", class)

	case "add-bib":
		_ = addBibCmd.Parse(os.Args[2:])
		if *addBibTitle == "" || *addBibAuthor == "" || *addBibType == "" || *addBibClass == 0 || *addBibYear == 0 {
			if *addBibTitle == "" {
				*addBibTitle = promptString("Title", true)
			}
			if *addBibAuthor == "" {
				*addBibAuthor = promptString("Author", true)
			}
			if *addBibPublisher == "" {
				*addBibPublisher = promptString("Publisher", false)
			}
			if *addBibType == "" {
				*addBibType = promptString("Type", true)
			}
			if *addBibClass == 0 {
				*addBibClass = promptInt("Classification Code Number", true)
			}
			if *addBibYear == 0 {
				*addBibYear = promptInt("Published Year", true)
			}
			if *addBibISBN == "" {
				*addBibISBN = promptString("ISBN", false)
			}
			if *addBibIndex == "" {
				*addBibIndex = promptString("BibIndex", false)
			}
		}
		// Construct date from year
		publishedDate := time.Date(*addBibYear, 1, 1, 0, 0, 0, 0, time.UTC)

		bib, err := app.BibService.AddBibliography(*addBibTitle, *addBibAuthor, *addBibPublisher, *addBibISBN, *addBibType, *addBibClass, publishedDate, *addBibTitleEn, *addBibAuthorEn, *addBibIndex)
		if err != nil {
			fmt.Printf("Error adding bibliography: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Bibliography added: %v\n", bib)

	case "add-review":
		_ = addReviewCmd.Parse(os.Args[2:])
		if *addReviewBibIndex == "" || *addReviewGoals == "" {
			if *addReviewBibIndex == "" {
				*addReviewBibIndex = promptString("BibIndex", true)
			}
			if *addReviewGoals == "" {
				*addReviewGoals = promptString("Goals", true)
			}
			if *addReviewSummary == "" {
				*addReviewSummary = promptString("Summary", false)
			}
		}

		// Resolve BibIndex to ID efficiently
		bib, err := app.BibService.FindByBibIndex(*addReviewBibIndex)
		if err != nil {
			fmt.Printf("Error finding bibliography with BibIndex %s: %v\n", *addReviewBibIndex, err)
			os.Exit(1)
		}
		if bib == nil {
			fmt.Printf("Bibliography with BibIndex %s not found\n", *addReviewBibIndex)
			os.Exit(1)
		}

		review, err := app.ReviewService.AddReview(bib.ID, *addReviewGoals, *addReviewSummary)
		if err != nil {
			fmt.Printf("Error adding review: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Review added: %v\n", review)

	case "update-review":
		_ = updateReviewCmd.Parse(os.Args[2:])
		if *updateReviewID == "" {
			*updateReviewID = promptString("Review UUID", true)
		}
		if *updateReviewGoals == "" && *updateReviewSummary == "" {
			// If neither is provided via flags, prompt for them
			if *updateReviewGoals == "" {
				*updateReviewGoals = promptString("New goals", false)
			}
			if *updateReviewSummary == "" {
				*updateReviewSummary = promptString("New summary", false)
			}
		}

		// Parse UUID
		reviewID, err := domain.ParseReviewID(*updateReviewID)
		if err != nil {
			fmt.Printf("Invalid review ID format: %v\n", err)
			os.Exit(1)
		}

		// Prepare optional fields
		var goals *string
		var summary *string
		if *updateReviewGoals != "" {
			goals = updateReviewGoals
		}
		if *updateReviewSummary != "" {
			summary = updateReviewSummary
		}

		// Validate at least one field is provided (after prompting)
		if goals == nil && summary == nil {
			fmt.Println("Please provide at least one field to update: -goals or -summary")
			// If interactive mode was used, they might have skipped both optional fields
			os.Exit(1)
		}

		review, err := app.ReviewService.UpdateReview(reviewID, goals, summary)
		if err != nil {
			fmt.Printf("Error updating review: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Review updated: %v\n", review)

	case "list":
		_ = listCmd.Parse(os.Args[2:])
		// Use user-specified limit and offset, or defaults
		bibs, err := app.BibService.ListBibliographies(*listLimit, *listOffset)
		if err != nil {
			fmt.Printf("Error listing bibliographies: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Bibliographies:")
		for _, b := range bibs {
			fmt.Printf("[%s] %s by %s (BibIndex: %s)\n", b.Type, b.Title, b.Author, b.BibIndex)
		}
		if len(bibs) == *listLimit && *listLimit > 0 {
			fmt.Printf("\nShowing %d items (use --limit and --offset to see more)\n", len(bibs))
		}

	default:
		fmt.Println("expected 'add-class', 'add-bib', 'add-review', 'update-review' or 'list' subcommands")
		os.Exit(1)
	}
}
