FROM golang:1.21.6 AS builder

ENV GO111MODULE=on

WORKDIR /app

RUN mkdir /var/log/entry/
COPY ../../ .
RUN --mount=type=cache,target=/go/pkg/mod/cache go mod download

RUN --mount=type=cache,target=/go/pkg/mod/cache CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o newsfeed_svc -mod=mod cmd/newsfeed/main.go

FROM builder AS test_env

RUN cp /app/configs/config.yml /app/config.yml
CMD ["/app/newsfeed_svc"]

# FROM builder AS live_env
# RUN cp /app/configs/config.yml /app/config.yml
# CMD ["/app/newsfeed_svc"]