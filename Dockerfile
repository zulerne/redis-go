FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o redis-server ./cmd/server

FROM gcr.io/distroless/static
COPY --from=build /app/redis-server /
EXPOSE 6379
ENTRYPOINT ["/redis-server"]
