FROM alpine:3.7 as application
WORKDIR /game-of-thrones

COPY ./dist/game-of-thrones/main /game-of-thrones
COPY ./cmd/config.docker.json /game-of-thrones/config.json
COPY ./migrations ./migrations

RUN apk update

EXPOSE 3033

ENTRYPOINT ./main