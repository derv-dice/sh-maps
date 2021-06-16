package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	. "sh-maps/config"
	"sh-maps/router"
	"sh-maps/view"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})

	var err error
	// Загрузка конфига
	if err = Cfg.Load(); err != nil {
		log.Fatalf("Не удалось загрузить конфиг. Ошибка: %v", err)
	}

	// Включение логирования
	if !Cfg.Log.NoLog {
		if err = LogToFile(Cfg.Log.LogDir); err != nil {
			log.Fatal(err)
		}
	}

	// Подготовка конфигов карт
	if err = PrepareBuildingsConfig(); err != nil {
		log.Fatalf("Не удалось подготовить конфиги карт. Ошибка: %v", err)
	}

	// Подготовка шаблонов и упаковка их в бинарник
	if err = view.PrepareTemplates(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	mux := router.Mux()

	go func() {
		log.Printf("Сервер запущен по адресу %s", Cfg.Server.Addr)
		log.Fatal(http.ListenAndServe(fmt.Sprintf("%s", Cfg.Server.Addr), mux))
	}()

	// Позаботимся о перехвате прерываний для корректной остановки сервиса
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGKILL, syscall.SIGSTOP, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGABRT)
	for {
		select {
		case <-stop:
			log.Println("Сервер остановлен")
			os.Exit(1)
		}
	}
}
