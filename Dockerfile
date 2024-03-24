FROM golang:1.21-bookworm AS build

WORKDIR /src

COPY go.mod ./

RUN go mod download

COPY . ./

RUN go build -o app .

FROM debian:bookworm

WORKDIR /

ARG port=8080

ENV PORT=$port

RUN apt update && apt install -y curl \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /src/app /app

EXPOSE $port

CMD [ "/app" ]
