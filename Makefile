antibf:
	cd scripts && ./genproto.sh
	cd cmd/antibf && go build -o antibf

cli:
	cd cmd/cli &&  go build -o cli

all: antibf cli

clean:
	rm -f cmd/antibf/antibf cmd/cli/cli
