FROM golang:1.21.6 AS builder

ENV GO111MODULE=on
WORKDIR /app

RUN mkdir /var/log/entry/
COPY ../../ .
RUN --mount=type=cache,target=/go/pkg/mod/cache go mod download

RUN --mount=type=cache,target=/go/pkg/mod/cache CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o web_app -mod=mod cmd/web_app/main.go


FROM builder AS test_env

RUN cp /app/configs/config.yml /app/config.yml
CMD ["/app/web_app"]

# FROM builder AS live_env
# RUN cp /app/configs/config.yml /app/config.yml
# CMD ["/app/web_app"]