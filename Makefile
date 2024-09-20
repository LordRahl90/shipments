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

dagger-build-env:
	dagger call build-env --source=. --verbose

dagger-test:
	dagger call test --source=. --verbose

bi: build-image
pi: push-image
dt: dagger-test