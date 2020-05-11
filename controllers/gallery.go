package controllers

import (
	"net/http"

	"lenslocked.com/views"
)

// Gallery defines the shape of a gallery
type Gallery struct {
	view *views.View
}

// NewGallery  returns a Gallery struct
func NewGallery() *Gallery {
	return &Gallery{
		view: views.NewView("bootstrap", "gallery/new"),
	}
}

// New method renders the new gallery page
func (g Gallery) New(w http.ResponseWriter, r *http.Request) {
	if err := g.view.Render(w, nil); err != nil {
		panic(err)
	}
}
