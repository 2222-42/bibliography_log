package main

import (
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
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	// Add Class Flags
	addClassCode := addClassCmd.Int("code", 0, "Classification Code Number (e.g. 56)")
	addClassName := addClassCmd.String("name", "", "Classification Name (e.g. Technology)")

	// Add Bib Flags
	addBibTitle := addBibCmd.String("title", "", "Title of the bibliography")
	addBibAuthor := addBibCmd.String("author", "", "Author of the bibliography")
	addBibType := addBibCmd.String("type", "", "Type (Book, Essay, Video, etc.)")
	addBibClass := addBibCmd.Int("class", 0, "Classification Code Number")
	addBibYear := addBibCmd.Int("year", 0, "Published Year (e.g. 2024)")
	addBibISBN := addBibCmd.String("isbn", "", "ISBN")
	addBibDesc := addBibCmd.String("desc", "", "Description")
	addBibTitleEn := addBibCmd.String("title-en", "", "English translation of title (required if title contains Japanese)")
	addBibAuthorEn := addBibCmd.String("author-en", "", "English translation of author (required if author contains Japanese)")
	addBibIndex := addBibCmd.String("bib-index", "", "Manual BibIndex (overrides auto-generation and bypasses English translation requirements)")

	// Add Review Flags
	addReviewBibIndex := addReviewCmd.String("bib-index", "", "BibIndex of the bibliography to review")
	addReviewGoals := addReviewCmd.String("goals", "", "Goals for reading (required)")
	addReviewSummary := addReviewCmd.String("summary", "", "Summary of the review")

	if len(os.Args) < 2 {
		fmt.Println("expected 'add-class', 'add-bib', 'add-review' or 'list' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add-class":
		addClassCmd.Parse(os.Args[2:])
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
		addBibCmd.Parse(os.Args[2:])
		if *addBibTitle == "" || *addBibAuthor == "" || *addBibType == "" || *addBibClass == 0 || *addBibYear == 0 {
			fmt.Println("Please provide required fields: -title, -author, -type, -class, -year")
			addBibCmd.PrintDefaults()
			os.Exit(1)
		}
		// Construct date from year
		publishedDate := time.Date(*addBibYear, 1, 1, 0, 0, 0, 0, time.UTC)

		bib, err := app.BibService.AddBibliography(*addBibTitle, *addBibAuthor, *addBibISBN, *addBibDesc, *addBibType, *addBibClass, publishedDate, *addBibTitleEn, *addBibAuthorEn, *addBibIndex)
		if err != nil {
			fmt.Printf("Error adding bibliography: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Bibliography added: %v\n", bib)

	case "add-review":
		addReviewCmd.Parse(os.Args[2:])
		if *addReviewBibIndex == "" || *addReviewGoals == "" {
			fmt.Println("Please provide required fields: -bib-index, -goals")
			addReviewCmd.PrintDefaults()
			os.Exit(1)
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

	case "list":
		listCmd.Parse(os.Args[2:])
		bibs, err := app.BibService.ListBibliographies()
		if err != nil {
			fmt.Printf("Error listing bibliographies: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Bibliographies:")
		for _, b := range bibs {
			fmt.Printf("[%s] %s by %s (BibIndex: %s)\n", b.Type, b.Title, b.Author, b.BibIndex)
		}

	default:
		fmt.Println("expected 'add-class', 'add-bib', 'add-review' or 'list' subcommands")
		os.Exit(1)
	}
}
