FROM golang:alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o shipments ./cmd/


FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /app/shipments shipments

EXPOSE 8080

ENTRYPOINT [ "./shipments" ]
