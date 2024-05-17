FROM golang:1.18-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY ./migrations /app/migrations

RUN go build -o main .

ENV ETHEREAL_EMAIL=ruthie.beier@ethereal.email
ENV ETHEREAL_PASSWORD=yWRW7aaB4dQMR8Sqsx

RUN go test -v ./...

FROM golang:1.18-alpine

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

CMD ["./main"]
