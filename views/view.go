package views

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"lenslocked.com/context"
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
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data{
			Yield: data,
		}
	}
	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	err := v.Template.ExecuteTemplate(&buf, v.Layout, vd)
	if err != nil {
		http.Error(w, "Something went wrong contact support for help", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
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
