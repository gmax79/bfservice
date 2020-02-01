abf:
	cd scripts && ./genproto.sh
	cd cmd/abf && go build -o abf

abfcli:
	cd cmd/abfcli && go build -o abfcli

cli: abfcli

all: abf abfcli

run: abf
	cd cmd/abf  && ./abf

clean:
	rm -f cmd/abg/abf cmd/abfcli/abfcli cmd/tests/tests

check:
	golangci-lint run --enable-all --disable wsl --disable lll --disable gochecknoglobals --disable gochecknoinits --disable gomnd

test:
	cd internal/netsupport && go test -v
	cd internal/buckets && go test -v