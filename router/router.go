package router

import (
	. "admin_template/config"
	"admin_template/handlers"
	"github.com/julienschmidt/httprouter"
)

func Mux() (mux *httprouter.Router) {
	mux = httprouter.New()
	mux.GET("/home", handlers.Home)

	mux.ServeFiles("/static/*filepath", Static)
	return
}
