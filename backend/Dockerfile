FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /api .

EXPOSE 3000

CMD ["./api"]