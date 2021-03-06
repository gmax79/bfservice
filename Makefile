.PHONY: abf abfcli all up check run stop test clean

abf:
	cd scripts && ./genproto.sh && cd ../cmd/abf && go build -o abf

abfcli:
	cd cmd/abfcli && go build -o ../../abfcli

abftests:
	cd cmd/tests && go build -o tests

all: abf abfcli abftests

up: abf
	cd cmd/abf  && ./abf

check:
	go vet ./...
	go fmt ./...
	golangci-lint run --enable-all --disable wsl --disable lll --disable gochecknoglobals --disable gochecknoinits --disable gomnd --disable interfacer

unittest:
	cd internal/netsupport && go test -v 2>&1
	cd internal/ratelimit && go test -v 2>&1
	cd cmd/tests && go test -v 2>&1
	cd internal/netsupport && go test -v -race 2>&1
	cd internal/ratelimit && go test -v -race 2>&1
	cd cmd/tests && go test -v -race 2>&1

run:
	docker-compose -f deployments/docker-compose.yml up -d

stop:
	docker-compose -f deployments/docker-compose.yml down

test:
	set -e; \
	export COMPOSE_IGNORE_ORPHANS=true; \
	exit_code=0; \
	docker-compose -f deployments/docker-compose.yml up -d; \
	docker-compose -f deployments/docker-compose.tests.yml up --exit-code-from tests || exit_code=$$?; \
	docker-compose -f deployments/docker-compose.yml -f deployments/docker-compose.tests.yml down; \
	make docker-clean; \
	echo "integration tests result: $$exit_code"; \
	exit $$exit_code;

docker:
	docker build -f build/package/abf.dockerfile -t gmax079/practice:abf .
	docker build -f build/package/tests.dockerfile -t gmax079/practice:abftests .

docker-push:
	docker push gmax079/practice:abf
	docker push gmax079/practice:abftests

docker-clean:
	docker rmi -f gmax079/practice:abftests gmax079/practice:abf

clean: docker-clean
	rm -f abfcli cmd/abf/abf cmd/tests/tests
