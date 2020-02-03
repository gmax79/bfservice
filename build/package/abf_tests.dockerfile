FROM golang:1.13.7 as builder
COPY . /app
WORKDIR /app/cmd/tests
RUN GOOS=linux go build -o tests .

FROM ubuntu:18.04
WORKDIR /root
COPY --from=builder /app/cmd/tests/tests ./
ENTRYPOINT "./tests"
