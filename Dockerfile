FROM golang:1.16-alpine AS build

WORKDIR /app

ENV CGO_ENABLED=0
ENV APIKEY_SENDGRID
ENV SENDER_EMAIL
ENV APIKEY_TELEGRAM_BOT
ENV TELEGRAM_CHANNEL_ID_VACCINE_ALERT
ENV TELEGRAM_CHANNEL_ID_VACCINE_ALERT_DEBUG

ADD go.mod go.sum . 

RUN go mod download

COPY . . 

ARG TARGETOS=linux
ARG TARGETARCH=arm64

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH}   go build -o bin/vaccine-alerts

# create volume
VOLUME [ "/app/shared" ]

# set entrypoint
ENTRYPOINT [ "./bin/vaccine-alerts" ]