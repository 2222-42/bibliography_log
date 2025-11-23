package main

import (
	"flag"
	"fmt"
	"os"
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
	addClassReq := &AddClassificationRequest{}
	addClassCmd.IntVar(&addClassReq.Code, "code", 0, "Classification Code Number (e.g. 56)")
	addClassCmd.StringVar(&addClassReq.Name, "name", "", "Classification Name (e.g. Technology)")

	// Add Bib Flags
	addBibReq := &AddBibliographyRequest{}
	addBibCmd.StringVar(&addBibReq.Title, "title", "", "Title of the bibliography")
	addBibCmd.StringVar(&addBibReq.Author, "author", "", "Author of the bibliography")
	addBibCmd.StringVar(&addBibReq.Publisher, "publisher", "", "Publisher of the bibliography")
	addBibCmd.StringVar(&addBibReq.Type, "type", "", "Type (Book, Essay, Video, etc.)")
	addBibCmd.IntVar(&addBibReq.ClassCode, "class", 0, "Classification Code Number")
	addBibCmd.IntVar(&addBibReq.Year, "year", 0, "Published Year (e.g. 2024)")
	addBibCmd.StringVar(&addBibReq.ISBN, "isbn", "", "ISBN")
	addBibCmd.StringVar(&addBibReq.TitleEn, "title-en", "", "English translation of title (required if title contains Japanese)")
	addBibCmd.StringVar(&addBibReq.AuthorEn, "author-en", "", "English translation of author (required if author contains Japanese)")
	addBibCmd.StringVar(&addBibReq.BibIndex, "bib-index", "", "Manual BibIndex (overrides auto-generation and bypasses English translation requirements)")

	// Add Review Flags
	addReviewReq := &AddReviewRequest{}
	addReviewCmd.StringVar(&addReviewReq.BibIndex, "bib-index", "", "BibIndex of the bibliography to review")
	addReviewCmd.StringVar(&addReviewReq.Goals, "goals", "", "Goals for reading (required)")
	addReviewCmd.StringVar(&addReviewReq.Summary, "summary", "", "Summary of the review")

	// Update Review Flags
	updateReviewReq := &UpdateReviewRequest{}
	updateReviewCmd.StringVar(&updateReviewReq.ReviewIDStr, "review-id", "", "UUID of the review to update (required)")
	updateReviewCmd.StringVar(&updateReviewReq.Goals, "goals", "", "New goals for reading (optional)")
	updateReviewCmd.StringVar(&updateReviewReq.Summary, "summary", "", "New summary of the review (optional)")

	// List Flags
	listReq := &ListBibliographiesRequest{}
	listCmd.IntVar(&listReq.Limit, "limit", 100, "Maximum number of items to display (default: 100, 0 for all)")
	listCmd.IntVar(&listReq.Offset, "offset", 0, "Number of items to skip (default: 0)")

	if len(os.Args) < 2 {
		fmt.Println("expected 'add-class', 'add-bib', 'add-review', 'update-review' or 'list' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add-class":
		_ = addClassCmd.Parse(os.Args[2:])
		addClassReq.PromptMissing()
		if err := addClassReq.Validate(); err != nil {
			fmt.Printf("Validation error: %v\n", err)
			addClassCmd.PrintDefaults()
			os.Exit(1)
		}

		class, err := app.BibService.AddClassification(addClassReq.Code, addClassReq.Name)
		if err != nil {
			fmt.Printf("Error adding classification: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Classification added: %v\n", class)

	case "add-bib":
		_ = addBibCmd.Parse(os.Args[2:])
		addBibReq.PromptMissing()
		if err := addBibReq.Validate(); err != nil {
			fmt.Printf("Validation error: %v\n", err)
			addBibCmd.PrintDefaults()
			os.Exit(1)
		}

		bib, err := app.BibService.AddBibliography(
			addBibReq.Title,
			addBibReq.Author,
			addBibReq.Publisher,
			addBibReq.ISBN,
			addBibReq.Type,
			addBibReq.ClassCode,
			addBibReq.ToPublishedDate(),
			addBibReq.TitleEn,
			addBibReq.AuthorEn,
			addBibReq.BibIndex,
		)
		if err != nil {
			fmt.Printf("Error adding bibliography: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Bibliography added: %v\n", bib)

	case "add-review":
		_ = addReviewCmd.Parse(os.Args[2:])
		addReviewReq.PromptMissing()
		if err := addReviewReq.Validate(); err != nil {
			fmt.Printf("Validation error: %v\n", err)
			addReviewCmd.PrintDefaults()
			os.Exit(1)
		}

		// Resolve BibIndex to ID efficiently
		bib, err := app.BibService.FindByBibIndex(addReviewReq.BibIndex)
		if err != nil {
			fmt.Printf("Error finding bibliography with BibIndex %s: %v\n", addReviewReq.BibIndex, err)
			os.Exit(1)
		}
		if bib == nil {
			fmt.Printf("Bibliography with BibIndex %s not found\n", addReviewReq.BibIndex)
			os.Exit(1)
		}

		review, err := app.ReviewService.AddReview(bib.ID, addReviewReq.Goals, addReviewReq.Summary)
		if err != nil {
			fmt.Printf("Error adding review: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Review added: %v\n", review)

	case "update-review":
		_ = updateReviewCmd.Parse(os.Args[2:])
		updateReviewReq.PromptMissing()
		if err := updateReviewReq.Validate(); err != nil {
			fmt.Printf("Validation error: %v\n", err)
			updateReviewCmd.PrintDefaults()
			os.Exit(1)
		}

		// Parse UUID
		reviewID, err := updateReviewReq.ParseID()
		if err != nil {
			fmt.Printf("Invalid review ID format: %v\n", err)
			os.Exit(1)
		}

		// Prepare optional fields
		var goals *string
		var summary *string
		if updateReviewReq.Goals != "" {
			goals = &updateReviewReq.Goals
		}
		if updateReviewReq.Summary != "" {
			summary = &updateReviewReq.Summary
		}

		review, err := app.ReviewService.UpdateReview(reviewID, goals, summary)
		if err != nil {
			fmt.Printf("Error updating review: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Review updated: %v\n", review)

	case "list":
		_ = listCmd.Parse(os.Args[2:])
		if err := listReq.Validate(); err != nil {
			fmt.Printf("Validation error: %v\n", err)
			listCmd.PrintDefaults()
			os.Exit(1)
		}

		bibs, err := app.BibService.ListBibliographies(listReq.Limit, listReq.Offset)
		if err != nil {
			fmt.Printf("Error listing bibliographies: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Bibliographies:")
		for _, b := range bibs {
			fmt.Printf("[%s] %s by %s (BibIndex: %s)\n", b.Type, b.Title, b.Author, b.BibIndex)
		}
		if len(bibs) == listReq.Limit && listReq.Limit > 0 {
			fmt.Printf("\nShowing %d items (use --limit and --offset to see more)\n", len(bibs))
		}

	default:
		fmt.Println("expected 'add-class', 'add-bib', 'add-review', 'update-review' or 'list' subcommands")
		os.Exit(1)
	}
}
