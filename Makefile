antibf:
	cd scripts && ./genproto.sh
	cd cmd/antibf && go build -o antibf

cli:
	cd cmd/cli &&  go build -o cli

all: antibf cli

run: antibf
	cd cmd/antibf  && ./antibf

clean:
	rm -f cmd/antibf/antibf cmd/cli/cli

check:
	golangci-lint run --enable-all
