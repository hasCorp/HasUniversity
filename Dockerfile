FROM golang:1.18-alpine3.16

ENV APP_HOME /go/src/hasuniversity

RUN mkdir -p "$APP_HOME"

WORKDIR "$APP_HOME"

COPY . "$APP_HOME"

EXPOSE 8200

CMD [ "go", "run", "main.go" ]