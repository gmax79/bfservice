abf:
	cd scripts && ./genproto.sh
	cd cmd/abf && go build -o abf

abfcli:
	cd cmd/abfcli && go build -o abfcli

cli: abfcli

all: abf abfcli

service: abf
	cd cmd/abf  && ./abf

clean:
	rm -f cmd/abg/abf cmd/abfcli/abfcli cmd/tests/tests

check:
	golangci-lint run --enable-all --disable wsl --disable lll --disable gochecknoglobals --disable gochecknoinits --disable gomnd

test:
	cd internal/netsupport && go test -v -race
	cd internal/buckets && go test -v -race
	cd cmd/tests && go test -v -race

docker:
	docker build -f build/package/abf.dockerfile -t abf . \
	&& docker tag abf gmax079/practice:abf

run:
	cd deployments && docker-compose -f docker-compose.yml up
