package handlers

import (
	"admin_template/view"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	tmplPageName = "home"
)

func Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	view.Render(w, tmplPageName, nil)
}
