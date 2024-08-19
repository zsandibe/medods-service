include .env
export

create-migrations:
	migrate create -ext sql -dir ./migrations -seq init_schema

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres


createdb: 
	docker exec -it postgres createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

dropdb:
	sudo docker exec -it postgres dropdb  $(DB_NAME)


migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST)/$(DB_NAME)?sslmode=disable" -verbose up 

migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST)/$(DB_NAME)?sslmode=disable" -verbose down



start:
	docker-compose up --build -d

stop:
	docker-compose down

all-containers-delete:
	docker ps -a -q | xargs -r docker rm -f


all-images-delete:
	docker images -q | xargs -r docker rmi -f


make test-v:
	go test -v ./...

make test-cover:
	go test -cover ./...

run:
	go run cmd/main.go


.PHONY: deps
deps:
	go mod tidy

.PHONY: swagger-init
swagger-init:
	@echo "Generate swagger gui"
	swag init -g  cmd/main.go



.PHONY: create-migrations postgres createdb dropdb migrate-up migrate-down start stop run 