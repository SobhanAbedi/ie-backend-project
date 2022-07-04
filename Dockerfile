# syntax=docker/dockerfile:1

FROM golang:1.18

WORKDIR /app
#ENV CGO_ENABLED=1

COPY go.mod ./
COPY go.sum ./
RUN go mod download

#RUN apk update
#RUN apk add gcc
#RUN apk add alpine-sdk
#RUN apk fix
#RUN apk update

COPY . ./

RUN go build -o /docker-mail-service

EXPOSE 8080

CMD [ "/docker-mail-service" ]