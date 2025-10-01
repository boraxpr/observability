FROM golang:1.25.1-alpine3.22 AS builder
ARG APP_NAME="otel"
WORKDIR /src/
COPY go.* ./
RUN go mod download
COPY ${APP_NAME}/main.go ./
RUN go build -o app .

FROM gcr.io/distroless/static-debian12:nonroot
ARG APP_PORT="8081"
WORKDIR /app/
COPY --from=builder /src/app ./
CMD ["./app"]
EXPOSE ${APP_PORT}
