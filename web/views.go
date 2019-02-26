package web

import (
  "html/template"
  "path/filepath"
  "net/http"
)

var LayoutDir string = "web/views/layouts"

func NewView(layout string, files ...string) *View {
  files = append(files, LayoutFiles()...)
  t, err := template.ParseFiles(files...)
  if err != nil {
    panic(err)
  }

  return &View{
    Template: t,
    Layout:   layout,
  }
}

type View struct {
  Template *template.Template
  Layout   string
}

func (v *View) Render(w http.ResponseWriter, data interface{}) error {
  return v.Template.ExecuteTemplate(w, v.Layout, data)
}


func LayoutFiles() []string {
  files, err := filepath.Glob(LayoutDir + "/*.html")
  if err != nil {
    panic(err)
  }
  return files
}

