.PHONY: bi
all: bi
	docker-compose up

start:
	@go run ./cmd/

test:
	@go test ./... --cover

build-image:
	@docker build -t gcr.io/neurons-be-test/shipments:latest .

push-image:
	@docker push gcr.io/neurons-be-test/shipments:latest

bi: build-image
pi: push-image