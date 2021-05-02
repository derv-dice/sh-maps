package view

import (
	"fmt"
	"net/http"
)

func Render(w http.ResponseWriter, tmplName string, data interface{}) (err error) {
	if TmplCompiled[tmplName] == nil {
		return fmt.Errorf("не удалось найти шаблон с указанным именем (%s)", tmplName)
	}
	return TmplCompiled[tmplName].Execute(w, data)
}
