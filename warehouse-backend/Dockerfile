FROM golang:1.23.4 AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o main .

FROM gcr.io/distroless/base

COPY --from=builder /app/main /app/main

CMD ["/app/main"]

