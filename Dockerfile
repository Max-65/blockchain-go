FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/node ./cmd/node

FROM alpine:3.23

WORKDIR /app
COPY --from=build /out/node /app/node

EXPOSE 3000 8080

ENTRYPOINT ["/app/node"]