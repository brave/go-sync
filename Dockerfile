FROM public.ecr.aws/bitnami/golang:1.20.10@sha256:024e66f81f70ddb560fba1e3e660dea4a76d7cef1dda76b27771241528694ac5 as builder

ARG VERSION
ARG BUILD_TIME
ARG COMMIT

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X github.com/brave/go-sync/server.version=${VERSION} -X github.com/brave/go-sync/server.buildTime=${BUILD_TIME} -X github.com/brave/go-sync/server.commit=${COMMIT}" \
    -o main .

FROM alpine:3.19@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b as artifact
RUN apk add --update ca-certificates # Certificates for SSL
COPY --from=builder /src/main main

EXPOSE 8295

CMD ["./main"]
