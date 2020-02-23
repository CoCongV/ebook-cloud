package render

import (
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin/render"
)

type (
	PongoProduction struct {
		Templates map[string]*pongo2.Template
		Path      string
	}

	Pongo struct {
		Template *pongo2.Template
		Name     string
		Data     interface{}
	}
)

var htmlContentType = []string{"application/html; charset=utf-8"}

func New(path string) *PongoProduction {
	pongo2.RegisterFilter("getFormat", getFormat)
	return &PongoProduction{map[string]*pongo2.Template{}, path}
}

func getFormat(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	ss := strings.Split(in.String(), ".")
	return pongo2.AsValue(ss[len(ss)-1]), nil
}

func (p PongoProduction) Instance(name string, data interface{}) render.Render {
	var t *pongo2.Template
	if tmpl, ok := p.Templates[name]; ok {
		t = tmpl
	} else {
		tmpl := pongo2.Must(pongo2.FromFile(path.Join(p.Path, name)))
		p.Templates[name] = tmpl
		t = tmpl
	}
	log.Println("test")
	log.Println(p.Templates)

	return Pongo{
		Template: t,
		Name:     name,
		Data:     data,
	}
}

func (p Pongo) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = htmlContentType
	}
}

func (p Pongo) Render(w http.ResponseWriter) error {
	ctx := pongo2.Context(p.Data.(pongo2.Context))
	return p.Template.ExecuteWriter(ctx, w)
}
