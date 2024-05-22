.PHONY: bi
all: bi
	docker-compose up

start:
	@go run ./cmd/

test:
	@go test ./... --cover

build-image:
	@docker build -t lordrahl/shipments:latest .

push-image:
	@docker push lordrahl/shipments:latest

bi: build-image
pi: push-image