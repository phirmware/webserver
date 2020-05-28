package views

import (
	"bytes"
	"html/template"
	"io"
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
func (v *View) Render(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	switch data.(type) {
	case Data:
	default:
		data = Data{
			Yield: data,
		}
	}
	var buf bytes.Buffer
	err := v.Template.ExecuteTemplate(&buf, v.Layout, data)
	if err != nil {
		http.Error(w, "Something went wrong contact support for help", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, nil)
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
