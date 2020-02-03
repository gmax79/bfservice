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

test: docker-tests
	cd internal/netsupport && go test -v -race
	cd internal/buckets && go test -v -race
	cd cmd/tests && go test -v -race
	docker run --rm tests
	docker rmi -f tests

docker:
	docker build -f build/package/abf.dockerfile -t abf . \
	&& docker tag abf gmax079/practice:abf

docker-tests:
	docker build -f build/package/abf_tests.dockerfile -t abftests .

docker-push:
	docker push gmax079/practice:abf

docker-clean:
	docker rmi -f tests abf gmax079/practice:abf

run:
	cd deployments && docker-compose -f docker-compose.yml up

clean: docker-clean
	rm -f cmd/abf/abf cmd/abfcli/abfcli cmd/tests/tests
