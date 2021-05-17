FROM golang:1.16-alpine AS build

RUN apk --no-cache add tzdata

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


FROM scratch AS final

COPY --from=build /app/bin /app/bin

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV TZ=Asia/Kolkata

# create volume
VOLUME [ "/app/shared" ]

# set entrypoint
ENTRYPOINT [ "/app/bin/vaccine-alerts" ]