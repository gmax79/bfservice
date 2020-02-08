abf:
	cd scripts && ./genproto.sh
	cd cmd/abf && go build -o abf

abfcli:
	cd cmd/abfcli && go build -o abfcli

all: abf abfcli

up: abf
	cd cmd/abf  && ./abf

check:
	golangci-lint run --enable-all --disable wsl --disable lll --disable gochecknoglobals --disable gochecknoinits --disable gomnd

unittest:
	cd internal/netsupport && go test -v -race
	cd internal/buckets && go test -v -race
	cd cmd/tests && go test -v -race

test: docker docker-tests unittest
	set -e; \
	export COMPOSE_IGNORE_ORPHANS=true; \
	docker-compose -f deployments/docker-compose.yml up -d; \
	exit_code=0; \
	docker run --network host --rm abftests || exit_code=$$?; \
	docker-compose -f deployments/docker-compose.yml down; \
	docker rmi -f abftests gmax079/practice:abf; \
	echo "integration_tests result: $$exit_code"; \
	exit $$exit_code

docker:
	docker build -f build/package/abf.dockerfile -t gmax079/practice:abf .

docker-tests:
	docker build -f build/package/tests.dockerfile -t abftests .

docker-push:
	docker push gmax079/practice:abf

docker-clean:
	docker rmi -f abftests gmax079/practice:abf

run:
	cd deployments && docker-compose -f docker-compose.yml up

clean: docker-clean
	rm -f cmd/abf/abf cmd/abfcli/abfcli cmd/tests/tests
