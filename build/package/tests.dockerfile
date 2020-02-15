FROM golang:1.13.7 as builder
COPY . /app
WORKDIR /app/cmd/tests
RUN GOOS=linux go build -o tests .
WORKDIR /app/cmd/abfcli
RUN GOOS=linux go build -o abfcli .

FROM ubuntu:18.04
WORKDIR /root
COPY build/package/tests_entrypoint.sh ./
COPY --from=builder /app/cmd/tests/tests ./
COPY --from=builder /app/cmd/abfcli/abfcli ./

ENTRYPOINT "./tests_entrypoint.sh"

