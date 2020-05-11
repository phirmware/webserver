package controllers

import (
	"lenslocked.com/views"
)

// Static defines the shape of our static struct
type Static struct {
	HomeView    *views.View
	ContactView *views.View
	FaqView     *views.View
}

// NewStatic returns the type static and parses static files
func NewStatic() *Static {
	return &Static{
		HomeView:    views.NewView("bootstrap", "static/home"),
		ContactView: views.NewView("bootstrap", "static/contact"),
		FaqView:     views.NewView("bootstrap", "static/faq"),
	}
}
