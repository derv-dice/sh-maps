package handlers

import (
	"fmt"
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"

	. "sh-maps/config"
	"sh-maps/view"
)

func AllMapsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	view.Render(w, view.PgMaps, NewBaseCtx().
		Add("Buildings", BuildingConfigs).
		Add("Description", UniversityDescription),
	)
}

func SingleMapHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	building := p.ByName("building")
	if building != "" && BuildingConfigs[building] == nil {
		http.NotFound(w, r)
		return
	}

	cfgUrl := fmt.Sprintf(RemoteAddrTmpl(), path.Join(Cfg.Server.Addr, BuildingCfgPath, building))
	view.Render(w, view.PgMap, NewBaseCtx().Add("CfgURL", cfgUrl))
}

// BuildingConfigHandler - Конфиг для определенного здания
func BuildingConfigHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("building")
	if name == "" || BuildingConfigs[name] == nil {
		http.Error(w, fmt.Sprintf("bad paremeter %q", "building"), http.StatusBadRequest)
		return
	}

	w.Write(BuildingConfigs[name].BuildingConfig)
}

// MapsIconsSVG - Костыль для получения иконок, требуемых на фронте для отрисовки кнопок на карте
func MapsIconsSVG(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	img, err := Static.Find("images/icons.svg")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(img)
}
