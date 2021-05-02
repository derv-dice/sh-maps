all: build

build:
	@echo "Запущена сборка проекта"
	@packr2 build
	@go build -o bin/server main.go
	@packr2 clean
	@rm admin_template
	@echo "Проект собран и находится в ./bin/server"

run: build
	@echo "Запуск сервера"
	@cd bin && ./server

test: build
	@echo "Запуск сервера в тестовом режиме без логирования в файл"
	@cd bin && ./server -nolog
