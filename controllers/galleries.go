package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"

	"lenslocked.com/views"
)

// Galleries defines the shape of a gallery
type Galleries struct {
	New *views.View
	gs  models.GalleryService
}

// GalleryForm defines the shape of the gallery form
type GalleryForm struct {
	Title string `schema:"title"`
}

// NewGalleries  returns a Galleries struct
func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  gs,
	}
}

// Create method creates a gallery resource
func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	gallery := models.Gallery{
		Title: form.Title,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	fmt.Fprintln(w, gallery)
}
