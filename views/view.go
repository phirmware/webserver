package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	// LayoutDir is the directory for layouts
	LayoutDir = "views/layouts/"
	// TemplateExt is our layout files extensions
	TemplateExt = ".html"
	// TemplateDir is our file directory
	TemplateDir = "views/"
)

// View defines the shape of the returned value from NewViews
type View struct {
	Template *template.Template
	Layout   string
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}

// Render is a utility function that renders our templtes
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

// NewView parses our templates and returns a pointer to View
func NewView(layout string, files ...string) *View {
	addTemplateExt(files)
	addTemplatePath(files)
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}
