abf:
	cd scripts && ./genproto.sh
	cd cmd/abf && go build -o abf

cli:
	cd cmd/cli &&  go build -o cli

all: abf cli

run: abf
	cd cmd/abf  && ./abf

clean:
	rm -f cmd/abg/abf cmd/cli/cli

check:
	golangci-lint run --enable-all
