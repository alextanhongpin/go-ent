include .env
export

entity:
	go run entgo.io/ent/cmd/ent init $(name)

up:
	@docker-compose up -d

down:
	@docker-compose down

gen:
	@go generate ./ent
