package handlers

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"

	. "sh-maps/config"
	"sh-maps/view"
)

func AllMapsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := view.Render(w, view.PgMaps, NewBaseCtx().
		Add("Buildings", BuildingConfigs).
		Add("Description", UniversityDescription),
	)

	if err != nil {
		log.Error(err)
	}
}

func SingleMapHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	building := p.ByName("building")
	if building != "" && BuildingConfigs[building] == nil {
		http.NotFound(w, r)
		return
	}

	cfgUrl := fmt.Sprintf("%s/%s", Cfg.Server.RemoteAddr, path.Join(BuildingCfgPath, building))
	err := view.Render(w, view.PgMap, NewBaseCtx().Add("CfgURL", cfgUrl))
	if err != nil {
		log.Error(err)
	}
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
