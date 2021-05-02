package view

import (
	. "admin_template/config"
	"html/template"
)

// Название страницы
const (
	pgHome = "home"
)

// Расположение шаблонов /templates/...
const (
	tmplBase       = "layout/base.gohtml"
	tmplNavigation = "layout/navigation.gohtml"

	tmplHomePage = "home/page.gohtml"
)

// tmplRequires - список шаблонов, требуемых для формирования конкретной страницы
var tmplRequires = map[string][]string{
	pgHome: {
		tmplBase, tmplNavigation, tmplHomePage,
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
