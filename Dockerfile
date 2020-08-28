FROM golang:1.15-alpine AS build-env

WORKDIR /app
COPY . .

ENV CGO_ENABLED=0

RUN go build \
    -ldflags "-X main.buildcommit=$(git rev-parse HEAD) -X main.buildtime=$(date +%Y%m%d.%H%M%S)" \
    -tags timetzdata \
    -o goapp \
    main.go


##################################################

FROM alpine:latest

WORKDIR /app

ARG APP_PORT=8000
ENV APP_PORT=${APP_PORT}

ARG DB_CONN="ktb@ktbserver:Passw0rd@tcp(ktbserver.mysql.database.azure.com)/thaichana?charset=utf8&parseTime=True&loc=Local"
ENV DB_CONN=${DB_CONN}

COPY --from=build-env /app/goapp ./goapp

EXPOSE ${APP_PORT}
ENTRYPOINT [ "./goapp" ]
