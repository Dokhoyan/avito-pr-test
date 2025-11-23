FROM golang:1.24.1-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/main ./cmd/http-server

FROM alpine:3.13 AS final

RUN apk --no-cache add ca-certificates

WORKDIR /

RUN mkdir -p /logs && chmod 777 /logs

COPY --from=build /bin/main /main

EXPOSE 8080

ENTRYPOINT ["/main"]