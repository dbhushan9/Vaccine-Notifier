FROM golang:1.16-alpine AS build

RUN apk --no-cache add tzdata

WORKDIR /app

ENV CGO_ENABLED=0

ADD go.mod go.sum . 

RUN go mod download

COPY . . 

ARG TARGETOS=linux
ARG TARGETARCH=arm64

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH}   go build -o bin/vaccine-alerts

RUN chmod 700 bin/vaccine-alerts

FROM scratch AS final

COPY --from=build /app/bin /app

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV TZ=Asia/Kolkata

ARG APIKEY_SENDGRID
ARG SENDER_EMAIL
ARG APIKEY_TELEGRAM_BOT
ARG TELEGRAM_CHANNEL_ID_VACCINE_ALERT
ARG TELEGRAM_CHANNEL_ID_VACCINE_ALERT_DEBUG

# create volume
VOLUME [ "/app/shared" ]

# set entrypoint
ENTRYPOINT [ "/app/vaccine-alerts" ]