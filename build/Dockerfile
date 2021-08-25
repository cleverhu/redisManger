FROM golang:1.18.0 as build

WORKDIR /app

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

COPY . .
RUN go mod tidy && go mod vendor

RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o redis-manger ./main.go

FROM docker.m.daocloud.io/alpine:3.15

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
WORKDIR /app

COPY --from=build /app/redis-manger  /app
COPY --from=build /app/config.yaml  /app

CMD ["./redis-manger"]
