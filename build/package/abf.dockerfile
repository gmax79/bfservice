FROM golang:1.13.7 as builder
COPY . /app
WORKDIR /app/cmd/abf
RUN GOOS=linux go build -o abf .

FROM ubuntu:18.04
RUN \
  apt-get update \
  && apt-get -y install gettext-base \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /root
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.2.1/wait ./wait
COPY build/package/abf_entrypoint.sh build/package/abf_config.json.template ./
COPY --from=builder /app/cmd/abf/abf ./
ENTRYPOINT "./abf_entrypoint.sh"
