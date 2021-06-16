all: build

build:
	@echo "Запущена сборка проекта"
	@packr2 build
	@go build -o bin/server main.go
	@packr2 clean
	@rm sh-maps
	@echo "Проект собран и находится в ./bin/server"

run: build
	@echo "Запуск сервера"
	@cd bin && ./server

docker_build: build
	@echo "Сборка docker образа"
	@sudo docker build -t sh-maps:v0.5.0 .

docker_run: docker_build
	@echo "Запуск контейнера"
	@sudo docker run --rm -it -p 127.0.0.1:8080:8080 -w /maps sh-maps:v0.5.0 ./server
