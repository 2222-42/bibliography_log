package main

import (
	"testing"
)

func TestAddClassificationRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request AddClassificationRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			request: AddClassificationRequest{Code: 1, Name: "Test"},
			wantErr: false,
		},
		{
			name:    "missing code",
			request: AddClassificationRequest{Code: 0, Name: "Test"},
			wantErr: true,
		},
		{
			name:    "missing name",
			request: AddClassificationRequest{Code: 1, Name: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.request.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("AddClassificationRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddBibliographyRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request AddBibliographyRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: AddBibliographyRequest{
				Title:     "Title",
				Author:    "Author",
				Type:      "Book",
				ClassCode: 1,
				Year:      2023,
			},
			wantErr: false,
		},
		{
			name: "missing title",
			request: AddBibliographyRequest{
				Author:    "Author",
				Type:      "Book",
				ClassCode: 1,
				Year:      2023,
			},
			wantErr: true,
		},
		{
			name: "missing author",
			request: AddBibliographyRequest{
				Title:     "Title",
				Type:      "Book",
				ClassCode: 1,
				Year:      2023,
			},
			wantErr: true,
		},
		{
			name: "missing type",
			request: AddBibliographyRequest{
				Title:     "Title",
				Author:    "Author",
				ClassCode: 1,
				Year:      2023,
			},
			wantErr: true,
		},
		{
			name: "missing class code",
			request: AddBibliographyRequest{
				Title:  "Title",
				Author: "Author",
				Type:   "Book",
				Year:   2023,
			},
			wantErr: true,
		},
		{
			name: "missing year",
			request: AddBibliographyRequest{
				Title:     "Title",
				Author:    "Author",
				Type:      "Book",
				ClassCode: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.request.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("AddBibliographyRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddReviewRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request AddReviewRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			request: AddReviewRequest{BibIndex: "B1", Goals: "Read"},
			wantErr: false,
		},
		{
			name:    "missing bib index",
			request: AddReviewRequest{Goals: "Read"},
			wantErr: true,
		},
		{
			name:    "missing goals",
			request: AddReviewRequest{BibIndex: "B1"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.request.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("AddReviewRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateReviewRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateReviewRequest
		wantErr bool
	}{
		{
			name:    "valid request with goals",
			request: UpdateReviewRequest{ReviewIDStr: "uuid", Goals: "New Goals"},
			wantErr: false,
		},
		{
			name:    "valid request with summary",
			request: UpdateReviewRequest{ReviewIDStr: "uuid", Summary: "New Summary"},
			wantErr: false,
		},
		{
			name:    "missing review id",
			request: UpdateReviewRequest{Goals: "New Goals"},
			wantErr: true,
		},
		{
			name:    "missing update fields",
			request: UpdateReviewRequest{ReviewIDStr: "uuid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.request.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("UpdateReviewRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListBibliographiesRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request ListBibliographiesRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			request: ListBibliographiesRequest{Limit: 10, Offset: 0},
			wantErr: false,
		},
		{
			name:    "negative limit",
			request: ListBibliographiesRequest{Limit: -1, Offset: 0},
			wantErr: true,
		},
		{
			name:    "negative offset",
			request: ListBibliographiesRequest{Limit: 10, Offset: -1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.request.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ListBibliographiesRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
