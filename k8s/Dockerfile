FROM golang:latest

WORKDIR /app
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY go.mod .
COPY go.sum .

RUN go mod download

ARG LOCATIONS_TABLE
ARG LOCATION_HISTORY_TABLE
ENV LOCATIONS_TABLE=$LOCATIONS_TABLE
ENV LOCATION_HISTORY_TABLE=$LOCATION_HISTORY_TABLE

COPY .aws/credentials /root/.aws/credentials
COPY .aws/config /root/.aws/config

ENTRYPOINT ["go", "run", "cmd/locsearch/main.go"]