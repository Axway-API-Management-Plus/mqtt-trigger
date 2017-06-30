FROM golang:alpine AS build
RUN apk add --no-cache make git
WORKDIR /app/src/mqtt-trigger
COPY Makefile .deps ./
RUN make deps-install

COPY . ./
RUN make

CMD [ "/app/src/mqtt-trigger/mqtt-trigger" ]

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=build /app/src/mqtt-trigger/mqtt-trigger /usr/bin
COPY ./mqtt-trigger.yml ./mqtt-trigger-test.yml ./

EXPOSE 1883
ENV PORT 1883
ENV MQTT_HOST 0.0.0.0
ENV MQTT_PORT 1883
ENV MQTT_USERNAME guest
ENV MQTT_PASSWORD guest

CMD ["mqtt-trigger"]
