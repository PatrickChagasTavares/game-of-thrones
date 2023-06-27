FROM golang:1.20 as builder
WORKDIR /build

COPY . .

RUN go mod download

RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS=-buildvcs=false go build -o ./main ./cmd/main.go


FROM alpine:3.7 as application
WORKDIR /game-of-thrones

COPY --from=builder /build/main .
COPY --from=builder /build/cmd/config.docker.json ./config.json
COPY --from=builder /build/migrations ./migrations

RUN apk update

EXPOSE 3033

CMD ["./main"]