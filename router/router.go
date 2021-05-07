package router

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
	"time"

	. "sh-maps/config"
	"sh-maps/router/handlers"
)

func Mux() http.Handler {
	r := httprouter.New()
	r.ServeFiles("/static/*filepath", Static)
	r.ServeFiles(fmt.Sprintf("/%s", filepath.Join(MapsStaticPath, "*filepath")), http.Dir("maps"))

	r.GET(fmt.Sprintf("/%s/:building", BuildingCfgPath), handlers.BuildingConfigHandler) // Получение конфига определенного корпуса
	r.GET("/maps", handlers.AllMapsHandler)                                              // Вывод списка всех корпусов
	r.GET("/maps/:building", handlers.SingleMapHandler)                                  // Отрисовка карты определенного корпуса
	r.GET("/mapplic-icons.svg", handlers.MapsIconsSVG)                       // Костыль для проброса файла с иконками для mapplic

	// Добавление middleware
	h := accessLog(r)
	h = recovery(h)
	return h
}

// recovery - Middleware, предотвращающий остановку приложения в случае критической ошибки
func recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				log.WithFields(log.Fields{
					"panic":  err,
					"method": r.Method,
					"ip":     r.RemoteAddr,
					"url":    r.URL.Path,
				}).Warning()
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// accessLog - Middleware, логирующий все входящие запросы
func accessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()  // Засекается момент времени, когда непосредственно началась обработка запроса
		next.ServeHTTP(w, r) // Обработка запроса

		// Логируем запрос только если это не загрузка статики.
		// Иначе запрос на отображение 1 страницы в логах займет 10+ строк
		if !Static.Has(r.URL.String()) {
			log.WithFields(log.Fields{
				"method": r.Method,
				"ip":     r.RemoteAddr,
				"url":    r.URL.String(),
				"time":   time.Since(start),
			}).Info()
		}
	})
}
