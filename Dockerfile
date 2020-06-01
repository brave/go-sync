FROM golang:1.14 as builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.6 as artifact
RUN apk add --update ca-certificates # Certificates for SSL
COPY --from=builder /src/main main

EXPOSE 8295

CMD ["./main"]
