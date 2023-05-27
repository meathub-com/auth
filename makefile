.DEFAULT_GOAL := swagger

install_swagger:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger:
	@echo Ensure you have the swagger CLI or this command will fail.
	@echo You can install the swagger CLI with: go get -u github.com/swaggo/swag/cmd/swag
	@echo ....

	swag init -g cmd/server.go

run-docker:
	docker-compose up
build-docker:
	docker-compose up --build