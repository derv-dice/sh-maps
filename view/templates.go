package view

import (
	"html/template"

	. "sh-maps/config"
)

// Название страницы
const (
	PgMaps = "maps"
	PgMap  = "map"
)

// Расположение шаблонов /templates/:tmpl_name.gohtml
const (
	tmplBase     = "layout/base.gohtml"
	tmplMapsPage = "pages/maps.gohtml"
	tmplMapPage  = "pages/map.gohtml"
)

// tmplRequires - список шаблонов, требуемых для формирования конкретной страницы
var tmplRequires = map[string][]string{
	PgMaps: {
		tmplBase, tmplMapsPage,
	},
	PgMap: {
		tmplBase, tmplMapPage,
	},
}

var TmplCompiled = map[string]*template.Template{}

func PrepareTemplates() (err error) {
	for k := range tmplRequires {
		t := template.New(k)

		// Собираем все требуемые шаблоны в 1 строку
		var ts string
		for n := range tmplRequires[k] {
			var s string
			if s, err = Templates.FindString(tmplRequires[k][n]); err != nil {
				return
			}
			ts += s
		}

		var tmpl *template.Template
		if tmpl, err = t.Parse(ts); err != nil {
			return
		}

		TmplCompiled[k] = tmpl
	}
	return
}
